package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	c "github.com/t1001001/prog-assig02/internal/constants"
)

var startTime = time.Now()

// HEAD request for RestCountries and OpenMeteo API
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

	resp, err := http.Head(statusURL)
	if err != nil {
		log.Printf("Error fetching %s: %v", statusURL, err)
		return "Unavailable"
	}
	defer resp.Body.Close()

	return resp.Status
}

// HEAD requests dont work on the Currency API, so we need a GET request here
func checkStatusGet(baseURL string) string {
	var statusURL string
	if strings.Contains(baseURL, c.CURRENCY_API_URL) {
		statusURL = baseURL + "/NOK"
		log.Printf("Checking the Currency API: %s", statusURL)
	} else {
		statusURL = baseURL
	}

	resp, err := http.Get(statusURL)
	if err != nil {
		log.Printf("Error fetching %s: %v", statusURL, err)
		return "Unavailable"
	}
	defer resp.Body.Close()

	return resp.Status
}

// StatusHandler handles the GET request for service status
func StatusHandler(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only the GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()

	// Check Firestore connection
	dbStatus := "500 Internal Server Error"
	_, err := client.Collection("webhooks").Limit(1).Documents(ctx).Next()
	if err == nil || err == iterator.Done {
		dbStatus = "200 OK"
	} else {
		log.Printf("Error checking Firestore connection: %v", err)
	}

	// Count registered webhooks
	webhookCount := 0
	iter := client.Collection("webhooks").Documents(ctx)
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error counting webhooks: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		webhookCount++
	}

	// Build status response
	status := c.Status{
		RestCountriesStatus: checkStatusHead(c.RESTCOUNTRIES_API_URL),
		OpenMeteoStatus:     checkStatusHead(c.OPENMETEO_API_URL),
		CurrencyStatus:      checkStatusGet(c.CURRENCY_API_URL),
		NotificationDB:      dbStatus,
		Webhooks:            webhookCount,
		Version:             strings.ToUpper(strings.Replace(c.VERSION, "/", "", -1)),
		Uptime:              int(time.Since(startTime).Seconds()),
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
