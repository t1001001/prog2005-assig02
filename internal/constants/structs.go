package constants

// Country represents the country information
type Country struct {
	Country       string   `json:"country" firestore:"country"`
	ISOCode       string   `json:"isoCode" firestore:"isoCode"`
	Features      Features `json:"features" firestore:"features"`
	LastRetrieval string   `json:"lastRetrieval" firestore:"lastRetrieval"`
}

// Features represents the "Features" field used in the "Country" struct
type Features struct {
	Temperature      float64            `json:"temperature" firestore:"temperature"`
	Precipitation    float64            `json:"precipitation" firestore:"precipitation"`
	Capital          string             `json:"capital" firestore:"capital"`
	Coordinates      Coordinates        `json:"coordinates" firestore:"coordinates"`
	Population       int                `json:"population" firestore:"population"`
	Area             float64            `json:"area" firestore:"area"`
	TargetCurrencies map[string]float64 `json:"targetCurrencies" firestore:"targetCurrencies"`
}

// Coordinates represents the "Coordinates" field used in the "Features" struct
type Coordinates struct {
	Latitude  float64 `json:"latitude" firestore:"latitude"`
	Longitude float64 `json:"longitude" firestore:"longitude"`
}

// Webhook represents the event notifications
type Webhook struct {
	ID      string `json:"id" firestore:"id"`
	URL     string `json:"url" firestore:"url"`
	Country string `json:"country" firestore:"country"`
	Event   string `json:"event" firestore:"event"`
}

// Status represents the status of the services
type Status struct {
	RestCountriesStatus string `json:"restcountries_status"`
	OpenMeteoStatus     string `json:"openmeteo_status"`
	CurrencyStatus      string `json:"currency_status"`
	NotificationDB      string `json:"notification_db"`
	Webhooks            int    `json:"webhooks"`
	Version             string `json:"version"`
	Uptime              int    `json:"uptime"`
}
