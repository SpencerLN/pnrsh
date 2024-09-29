package pnr

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
)

const (
	reqBody      = `{"recLoc":"{conf}","lName":"{lname}"}`
	reqUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36'"
	reqEndpoint  = "https://apis.alaskaair.com/guestServices/flightChanges/bff/ChangeReservation/ReservationDetails"
)

var (
	client = http.Client{}
)

func generateReqBody(lastName, confirmationCode string) string {
	body := reqBody
	body = strings.Replace(body, "{conf}", confirmationCode, -1)
	body = strings.Replace(body, "{lname}", lastName, -1)
	return body
}

func buildRequest(endpoint, body string) (*http.Request, io.ReadSeeker) {
	reader := strings.NewReader(body)
	return buildRequestWithBodyReader(endpoint, reader)
}

func buildRequestWithBodyReader(endpoint string, body io.Reader) (*http.Request, io.ReadSeeker) {
	var bodyLen int

	type lenner interface {
		Len() int
	}
	if lr, ok := body.(lenner); ok {
		bodyLen = lr.Len()
	}

	req, _ := http.NewRequest("POST", endpoint, body)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Api-Version", "2")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://www.alaskaair.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://www.alaskaair.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", reqUserAgent)
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	if bodyLen > 0 {
		req.Header.Set("Content-Length", strconv.Itoa(bodyLen))
	}

	var seeker io.ReadSeeker
	if sr, ok := body.(io.ReadSeeker); ok {
		seeker = sr
	} else {
		seeker = aws.ReadSeekCloser(body)
	}

	return req, seeker
}

func sendRequest(lastName, confirmationCode string) ([]byte, error) {
	req, _ := buildRequest(reqEndpoint, generateReqBody(lastName, confirmationCode))

	res, _ := client.Do(req)

	if res.StatusCode != 200 {
		return []byte{}, errors.New("status code was not 200")
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func performRequest(lastName, confirmationCode string) (res AlaskaManagePnrResponse, err error) {

	data, err := sendRequest(lastName, confirmationCode)
	if err != nil {
		return res, err
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return res, err
	}

	// // Open our jsonFile
	// jsonFile, err := os.Open("/Users/surehit/Desktop/as-pnr.json")
	// // if we os.Open returns an error then handle it
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Successfully Opened users.json")
	// // defer the closing of our jsonFile so that we can parse it later on
	// defer jsonFile.Close()
	// // read our opened jsonFile as a byte array.
	// byteValue, _ := ioutil.ReadAll(jsonFile)

	// // we initialize our Users array
	// var pnrData AlaskaManagePnrResponse

	// // we unmarshal our byteArray which contains our
	// // jsonFile's content into 'users' which we defined above
	// json.Unmarshal(byteValue, &pnrData)

	return res, nil
}

func convertResponse(res AlaskaManagePnrResponse) (pnr PNR) {

	// convertRemarks(res, &pnr)
	convertFlights(res, &pnr)
	convertOSIs(res, &pnr)
	convertPassengers(res, &pnr)
	convertTickets(res, &pnr)
	convertItinerary(res, &pnr)

	return pnr
}
