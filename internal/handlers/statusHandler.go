package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	c "github.com/t1001001/prog-assig02/internal/constants"
)

var startTime = time.Now()

// send HEAD to the RestCountries and OpenMeteo API
func checkStatusHead(baseURL string) string {
	var statusURL string
	if strings.Contains(baseURL, c.RESTCOUNTRIES_API_URL) {
		statusURL = baseURL + "/all"
		log.Printf("Checking the RestCountries API: %s", statusURL)
	} else if strings.Contains(baseURL, c.OPENMETEO_API_URL) {
		statusURL = baseURL + "/forecast"
		log.Printf("Checking the Open-Meteo API: %s", statusURL)
	} else {
		statusURL = baseURL
	}

	// Check if the API is available
	resp, err := http.Head(statusURL)
	if err != nil {
		log.Printf("Error fetching %s: %v", statusURL, err)
		return "Unavailable"
	}
	defer resp.Body.Close()

	return resp.Status
}

// Send GET to the Currency API
func checkStatusGet(baseURL string) string {
	var statusURL string
	if strings.Contains(baseURL, c.CURRENCY_API_URL) {
		statusURL = baseURL + "/NOK"
		log.Printf("Checking the Currency API: %s", statusURL)
	} else {
		statusURL = baseURL
	}

	// Check if the API is available
	resp, err := http.Get(statusURL)
	if err != nil {
		log.Printf("Error fetching %s: %v", statusURL, err)
		return "Unavailable"
	}
	defer resp.Body.Close()

	return resp.Status
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Only the GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// status response
	status := c.Status{
		RestCountriesStatus: checkStatusHead(c.RESTCOUNTRIES_API_URL),
		OpenMeteoStatus:     checkStatusHead(c.OPENMETEO_API_URL),
		CurrencyStatus:      checkStatusGet(c.CURRENCY_API_URL),
		NotificationDB:      "",
		Webhooks:            0,
		Version:             c.VERSION,
		Uptime:              int(time.Since(startTime).Seconds()),
	}

	// Setting header
	w.Header().Set("Content-Type", "application/json")

	// Encode the response
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
