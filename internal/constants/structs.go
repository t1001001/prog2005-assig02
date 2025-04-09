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
	Temperature      bool     `json:"temperature" firestore:"temperature"`
	Precipitation    bool     `json:"precipitation" firestore:"precipitation"`
	Capital          bool     `json:"capital" firestore:"capital"`
	Coordinates      bool     `json:"coordinates" firestore:"coordinates"`
	Population       bool     `json:"population" firestore:"population"`
	Area             bool     `json:"area" firestore:"area"`
	TargetCurrencies []string `json:"targetCurrencies" firestore:"targetCurrencies"`
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
	Time    string `json:"time" firestore:"time"`
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

// RegistrationResponse represents the response format for a registration
type RegistrationResponse struct {
	ID         string   `json:"id"`
	Country    string   `json:"country"`
	ISOCode    string   `json:"isoCode"`
	Features   Features `json:"features"`
	LastChange string   `json:"lastChange"`
}

// DashboardResponse is returned by GET /dashboards/{id}
type DashboardResponse struct {
	Country       string           `json:"country"`
	ISOCode       string           `json:"isoCode"`
	Features      EnrichedFeatures `json:"features"`
	LastRetrieval string           `json:"lastRetrieval"`
}

// EnrichedFeatures contains fetched external data
type EnrichedFeatures struct {
	Temperature      float64            `json:"temperature,omitempty"`
	Precipitation    float64            `json:"precipitation,omitempty"`
	Capital          string             `json:"capital,omitempty"`
	Coordinates      Coordinates        `json:"coordinates,omitempty"`
	Population       int                `json:"population,omitempty"`
	Area             float64            `json:"area,omitempty"`
	TargetCurrencies map[string]float64 `json:"targetCurrencies,omitempty"`
}
