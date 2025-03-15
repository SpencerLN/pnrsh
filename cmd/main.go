package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	aeromexicoEnabled = true
	aircanadaEnabled  = true
	deltaEnabled      = true
	unitedEnabled     = true
	virginEnabled     = true
	alaskaEnabled     = true

	// Session constants
	sessionName    = "pnrsh-session"
	sessionUserKey = "user-email"
)

var (
	commitHash = os.Getenv("HEROKU_SLUG_COMMIT")
	// OAuth2 configuration
	googleOAuthConfig *oauth2.Config
	// Session store for cookies
	sessionStore *sessions.CookieStore
	// Allowed email addresses
	allowedEmails map[string]bool
)

func initOAuth() {
	// Initialize allowed emails from environment variable
	allowedEmails = make(map[string]bool)
	emailList := os.Getenv("ALLOWED_EMAILS")
	if emailList != "" {
		for _, email := range strings.Split(emailList, ",") {
			allowedEmails[strings.TrimSpace(email)] = true
		}
	}

	// Get OAuth credentials from environment
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Println("OAuth credentials not set, authentication will not work")
	}

	// Random key for cookie store
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		sessionKey = "pnrsh-random-key"
		log.Println("Warning: Using default session key. Set SESSION_KEY environment variable in production.")
	}
	sessionStore = sessions.NewCookieStore([]byte(sessionKey))

	// Create OAuth config
	redirectURL := os.Getenv("OAUTH_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = "http://" + listenAddress() + "/auth/callback"
	}

	googleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func listenAddress() string {
	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "127.0.0.1"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return listenAddress + ":" + port
}

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for login-related paths
		if r.URL.Path == "/auth/login" || r.URL.Path == "/auth/callback" {
			next.ServeHTTP(w, r)
			return
		}

		session, _ := sessionStore.Get(r, sessionName)
		email, ok := session.Values[sessionUserKey].(string)
		if !ok || email == "" {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		// Check if the email is in the allowed list
		if !allowedEmails[email] {
			http.Error(w, "Unauthorized: Email not in allowed list", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoginHandler initiates the OAuth flow
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// CallbackHandler handles the OAuth callback
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Simple struct to parse the email from response
	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the email is in the allowed list
	if !allowedEmails[userInfo.Email] {
		http.Error(w, "Unauthorized: Email not in allowed list", http.StatusUnauthorized)
		return
	}

	// Save the email in session
	session, _ := sessionStore.Get(r, sessionName)
	session.Values[sessionUserKey] = userInfo.Email
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutHandler logs the user out
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionName)
	session.Values[sessionUserKey] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	// Initialize OAuth
	initOAuth()
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/help", HelpHandler).Methods("GET")

	// Add auth routes
	r.HandleFunc("/auth/login", LoginHandler).Methods("GET")
	r.HandleFunc("/auth/callback", CallbackHandler).Methods("GET")
	r.HandleFunc("/auth/logout", LogoutHandler).Methods("GET")
	if aeromexicoEnabled {
		r.HandleFunc("/aeromexico", AeromexicoHomeHandler).Methods("GET")
		r.HandleFunc("/aeromexico", AeromexicoRetrieveHandler).Methods("POST")
	}

	if aircanadaEnabled {
		r.HandleFunc("/aircanada", AircanadaHomeHandler).Methods("GET")
		r.HandleFunc("/aircanada", AircanadaRetrieveHandler).Methods("POST")
	}

	if deltaEnabled {
		r.HandleFunc("/delta", DeltaHomeHandler).Methods("GET")
		r.HandleFunc("/delta", DeltaRetrieveHandler).Methods("POST")
	}

	if unitedEnabled {
		r.HandleFunc("/united", UnitedHomeHandler).Methods("GET")
		r.HandleFunc("/united", UnitedRetrieveHandler).Methods("POST")
	}

	if virginEnabled {
		r.HandleFunc("/virgin", VirginHomeHandler).Methods("GET")
		r.HandleFunc("/virgin", VirginRetrieveHandler).Methods("POST")
	}

	if alaskaEnabled {
		r.HandleFunc("/alaska", AlaskaHomeHandler).Methods("GET")
		r.HandleFunc("/alaska", AlaskaRetrieveHandler).Methods("POST")
	}
	r.Use(AuthMiddleware)

	srv := &http.Server{
		Handler: r,
		Addr:    listenAddress(),

		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Println("Visit", "http://"+listenAddress(), "to use the app!")
	log.Fatal(srv.ListenAndServe())
}
