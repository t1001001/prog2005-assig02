package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	c "github.com/t1001001/prog-assig02/internal/constants"
)

// DashboardsHandler handles the GET request for the /dashboards/{id} endpoint
func DashboardsHandler(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()
	id := strings.TrimPrefix(r.URL.Path, c.ROOT+c.VERSION+c.DASHBOARDS_PATH)
	id = strings.Trim(id, "/")

	if id == "" {
		http.Error(w, "Missing ID in path", http.StatusBadRequest)
		return
	}

	// Retrieve the corresponding configuration from Firestore
	doc, err := client.Collection("registrations").Doc(id).Get(ctx)
	if err != nil {
		http.Error(w, "Configuration not found", http.StatusNotFound)
		return
	}

	var config c.Country
	if err := doc.DataTo(&config); err != nil {
		http.Error(w, "Failed to decode configuration", http.StatusInternalServerError)
		return
	}

	log.Printf("[DASHBOARD] Loaded config for: %s (%s)", config.Country, config.ISOCode)

	// Initialize the dashboard response structure
	response := c.DashboardResponse{
		Country:       config.Country,
		ISOCode:       config.ISOCode,
		LastRetrieval: config.LastRetrieval,
		Features:      c.EnrichedFeatures{},
	}

	// Respond with the formatted dashboard data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	out, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Write(out)
}
