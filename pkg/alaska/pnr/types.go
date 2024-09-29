package pnr

import "time"

type MyTime struct {
	time.Time
}

func (self *MyTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// Get rid of the quotes "" around the value.
	// A second option would be to include them
	// in the date format string instead, like so below:
	//   time.Parse(`"`+time.RFC3339Nano+`"`, s)
	s = s[1 : len(s)-1]
	if s == "0001-01-01T00:00:00" {
		s = "0001-01-01T00:00:00Z"
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.999999999Z0700", s)
	}
	self.Time = t
	return
}

type PNR struct {
	Remarks    []Remark
	Flights    []Flight
	Passengers []Passenger
	Tickets    []Ticket
	Itinerary  Itinerary
	SMCalcLink string
	OSIs       []OSI
}

type Itinerary struct {
	Type           string
	Origin         string
	Destination    string
	MatchesTickets bool
}

type Remark struct {
	FreeFormText string
	RemarkType   string
}

type Flight struct {
	OriginAirportCode      string
	DestinationAirportCode string
	Distance               int
	CurrentActionCode      string
	PreviousActionCode     string
	Status                 string
	MarketingAirlineCode   string
	OperatingAirlineCode   string
	ClassOfService         string
	Cabin                  string
	UpgradeStatus          string
	FUpgradeStatus         string
	PcUpgradeStatus        string
	ScheduledDeparture     string
	ScheduledArrival       string
	FlightNumber           string
	IsDisrupted            bool
	IsFlown                bool
	SSRs                   []SSR

	// On AM, this is a per-segment field... dunno.
	FareBasis string
	FareName  string
}

type Passenger struct {
	Name                   string
	Status                 string
	OverbookingEligible    bool
	SkyPriority            bool
	TierStatus             string
	BenefitCodes           string
	FfNumber               string
	FeeForSeatSelection    bool
	FeeForAmPlusUpgrade    bool
	FeeForPreferredUpgrade bool
	SSRs                   []SSR
	OSIs                   []OSI
}

type Coupon struct {
	Index             int    `json:"index"`
	DepartureDateTime string `json:"departureDateTime"`
	Origin            string `json:"origin"`
	Destination       string `json:"destination"`
	Status            string `json:"status"`
}

type SSR struct {
	AirlineCode string
	Type        string
	Remark      string
	Status      string
	FlightNum   string
	FlightDate  string
	Id          string
}

type OSI struct {
	FreeText   string
	FullText   string
	VendorCode string
	Id         string
}

type Ticket struct {
	Number                string
	CouponNumber          string
	Status                bool
	PreviousStatus        string
	PassengerName         string
	NumCoupons            uint64
	Coupons               []Coupon
	RelatedDocumentNumber string
	OriginDestination     string
	TotalCost             string
	Type                  string
	Designator            string
}

type AlaskaManagePnrResponse struct {
	Itinerary struct {
		TripType        string      `json:"tripType"`
		Origin          string      `json:"origin"`
		Destination     interface{} `json:"destination"`
		ItinerarySlices []struct {
			Origin      string `json:"origin"`
			Destination string `json:"destination"`
			Segments    []struct {
				DepartureAirport              string      `json:"departureAirport"`
				DepartureCity                 string      `json:"departureCity"`
				ArrivalAirport                string      `json:"arrivalAirport"`
				ArrivalCity                   string      `json:"arrivalCity"`
				OperatingAirlineCode          string      `json:"operatingAirlineCode"`
				OperatingAirlineName          string      `json:"operatingAirlineName"`
				OperatingFlightNumber         string      `json:"operatingFlightNumber"`
				ActionCode                    string      `json:"actionCode"`
				DisclosureText                string      `json:"disclosureText"`
				ScheduledDepartureDateTime    string      `json:"scheduledDepartureDateTime"`
				ScheduledDepartureDateTimeUtc MyTime      `json:"scheduledDepartureDateTimeUtc"`
				EstimatedDepartureDateTime    interface{} `json:"estimatedDepartureDateTime"`
				EstimatedDepartureDateTimeUtc interface{} `json:"estimatedDepartureDateTimeUtc"`
				ScheduledArrivalDateTime      string      `json:"scheduledArrivalDateTime"`
				ScheduledArrivalDateTimeUtc   MyTime      `json:"scheduledArrivalDateTimeUtc"`
				EstimatedArrivalDateTime      interface{} `json:"estimatedArrivalDateTime"`
				EstimatedArrivalDateTimeUtc   interface{} `json:"estimatedArrivalDateTimeUtc"`
				ScheduledDurationInMinutes    int         `json:"scheduledDurationInMinutes"`
				EstimatedDurationInMinutes    int         `json:"estimatedDurationInMinutes"`
				Status                        interface{} `json:"status"`
				Cabin                         string      `json:"cabin"`
				EquipmentName                 string      `json:"equipmentName"`
				Distance                      int         `json:"distance"`
				ClassOfService                string      `json:"classOfService"`
				IsArnk                        bool        `json:"isArnk"`
				SpecialServiceRequests        []struct {
					ID                 string      `json:"id"`
					ServiceCode        string      `json:"serviceCode"`
					FreeText           string      `json:"freeText"`
					ServiceDescription interface{} `json:"serviceDescription"`
					ActionCode         string      `json:"actionCode"`
					VendorCode         string      `json:"vendorCode"`
					FlightNumber       string      `json:"flightNumber"`
					FlightDate         string      `json:"flightDate"`
					Origin             string      `json:"origin"`
					Destination        string      `json:"destination"`
				} `json:"specialServiceRequests"`
				Sequence                       int  `json:"sequence"`
				IsDisrupted                    bool `json:"isDisrupted"`
				IsFlown                        bool `json:"isFlown"`
				IsWaitingUpgradeFirstClass     bool `json:"isWaitingUpgradeFirstClass"`
				IsWaitingUpgradePremiumClass   bool `json:"isWaitingUpgradePremiumClass"`
				IsConfirmedUpgradeFirstClass   bool `json:"isConfirmedUpgradeFirstClass"`
				IsConfirmedUpgradePremiumClass bool `json:"isConfirmedUpgradePremiumClass"`
			} `json:"segments"`
			PreviousItinerary   []interface{} `json:"previousItinerary"`
			HistoricItineraries []interface{} `json:"historicItineraries"`
			IsInternationalTrip bool          `json:"isInternationalTrip"`
			IsBagCheckedIn      bool          `json:"isBagCheckedIn"`
			IsBagTagCreated     bool          `json:"isBagTagCreated"`
			IsOAInitiatedTrip   bool          `json:"isOAInitiatedTrip"`
			DisruptionType      string        `json:"disruptionType"`
		} `json:"itinerarySlices"`
		MatchesTickets bool `json:"matchesTickets"`
	} `json:"itinerary"`
	Passengers []struct {
		FirstName     string `json:"firstName"`
		LastName      string `json:"lastName"`
		NameRefNumber string `json:"nameRefNumber"`
		TierStatus    string `json:"tierStatus"`
		LoyaltyNumber string `json:"loyaltyNumber"`
		Tickets       []struct {
			Number     string `json:"number"`
			Type       string `json:"type"`
			Designator string `json:"designator"`
			Coupons    []struct {
				Index             int    `json:"index"`
				DepartureDateTime string `json:"departureDateTime"`
				Origin            string `json:"origin"`
				Destination       string `json:"destination"`
				Status            string `json:"status"`
			} `json:"coupons"`
			Payments []struct {
				Type        string `json:"type"`
				Certificate string `json:"certificate"`
			} `json:"payments"`
			IsActive bool `json:"isActive"`
		} `json:"tickets"`
		SeatNumbers            [][]interface{} `json:"seatNumbers"`
		SpecialServiceRequests []struct {
			ID                 string      `json:"id"`
			ServiceCode        string      `json:"serviceCode"`
			FreeText           string      `json:"freeText"`
			ServiceDescription interface{} `json:"serviceDescription"`
			ActionCode         string      `json:"actionCode"`
			VendorCode         string      `json:"vendorCode"`
			FlightNumber       string      `json:"flightNumber"`
			FlightDate         string      `json:"flightDate"`
			Origin             string      `json:"origin"`
			Destination        string      `json:"destination"`
		} `json:"specialServiceRequests"`
		OSIs []struct {
			ID         string `json:"id"`
			FreeText   string `json:"freeText"`
			FullText   string `json:"fullText"`
			VendorCode string `json:"vendorCode"`
		} `json:"osis"`
		TicketsMatchesItinerary bool `json:"ticketsMatchesItinerary"`
	} `json:"passengers"`
	OSIs []struct {
		ID         string `json:"id"`
		FreeText   string `json:"freeText"`
		FullText   string `json:"fullText"`
		VendorCode string `json:"vendorCode"`
	} `json:"osis"`
	IsGroupBooking bool `json:"isGroupBooking"`
	Tickets        []struct {
		Number     string `json:"number"`
		Type       string `json:"type"`
		Designator string `json:"designator"`
		Coupons    []struct {
			Index             int    `json:"index"`
			DepartureDateTime string `json:"departureDateTime"`
			Origin            string `json:"origin"`
			Destination       string `json:"destination"`
			Status            string `json:"status"`
		} `json:"coupons"`
		Payments []struct {
			Type        string `json:"type"`
			Certificate string `json:"certificate"`
		} `json:"payments"`
		IsActive bool `json:"isActive"`
	} `json:"tickets"`
	IsUpgradeEligible               bool `json:"isUpgradeEligible"`
	IsOptedInForFirstClassUpgrade   bool `json:"isOptedInForFirstClassUpgrade"`
	IsOptedInForPremiumClassUpgrade bool `json:"isOptedInForPremiumClassUpgrade"`
}
