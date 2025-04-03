package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"

	c "github.com/t1001001/prog-assig02/internal/constants"
)

// getCurrentTimestamp returns the current time formatted as "YYYYMMDD HH:MM"
func getCurrentTimestamp() string {
	return time.Now().Format("20060102 15:04")
}

// RegistrationsHandler handles all HTTP methods for the /registrations endpoint
func RegistrationsHandler(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	fmt.Println(">>> Received request method:", r.Method, "at", r.URL.Path)
	ctx := context.Background()
	path := strings.TrimPrefix(r.URL.Path, c.ROOT+c.VERSION+c.REGISTRATIONS_PATH)
	id := strings.Trim(path, "/")

	switch r.Method {
	case http.MethodPost:
		// Handle POST request to register a new configuration
		var config c.Country
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		config.LastRetrieval = getCurrentTimestamp()
		id := uuid.New().String()
		_, err := client.Collection("registrations").Doc(id).Set(ctx, config)
		if err != nil {
			http.Error(w, "Failed to store configuration", http.StatusInternalServerError)
			return
		}
		fmt.Println("Successfully registered config for:", config.Country)
		resp := map[string]string{"id": id, "lastChange": config.LastRetrieval}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	case http.MethodGet:
		if id != "" {
			// Handle GET request for a specific registration
			doc, err := client.Collection("registrations").Doc(id).Get(ctx)
			if err != nil {
				http.Error(w, "Configuration not found", http.StatusNotFound)
				return
			}
			var config c.Country
			e := doc.DataTo(&config)
			if e != nil {
				fmt.Println("Error mapping document:", e)
			} else {
				fmt.Println("Loaded config for:", config.Country)
			}
			response := map[string]interface{}{
				"id":         doc.Ref.ID,
				"country":    config.Country,
				"isoCode":    config.ISOCode,
				"features":   config.Features,
				"lastChange": config.LastRetrieval,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			// Handle GET request to retrieve all registrations
			fmt.Println("Handling GET request for all registrations")
			docs, err := client.Collection("registrations").Documents(ctx).GetAll()
			if err != nil {
				http.Error(w, "Failed to fetch configurations", http.StatusInternalServerError)
				return
			}
			var configs []map[string]interface{}
			for _, doc := range docs {
				var config c.Country
				doc.DataTo(&config)
				entry := map[string]interface{}{
					"id":         doc.Ref.ID,
					"country":    config.Country,
					"isoCode":    config.ISOCode,
					"features":   config.Features,
					"lastChange": config.LastRetrieval,
				}
				configs = append(configs, entry)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(configs)
		}

	case http.MethodPut:
		// Handle PUT request to replace a specific registration
		if id == "" {
			http.Error(w, "ID required", http.StatusBadRequest)
			return
		}
		var config c.Country
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		config.LastRetrieval = getCurrentTimestamp()
		_, err := client.Collection("registrations").Doc(id).Set(ctx, config)
		if err != nil {
			http.Error(w, "Failed to update configuration", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodDelete:
		// Handle DELETE request to remove a specific registration
		if id == "" {
			http.Error(w, "ID required", http.StatusBadRequest)
			return
		}
		_, err := client.Collection("registrations").Doc(id).Delete(ctx)
		if err != nil {
			http.Error(w, "Failed to delete configuration", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodHead:
		// Handle HEAD request to check existence of a registration
		if id != "" {
			_, err := client.Collection("registrations").Doc(id).Get(ctx)
			if err != nil {
				http.Error(w, "Configuration not found", http.StatusNotFound)
				return
			}
		} else {
			// Just return OK for the collection
			w.WriteHeader(http.StatusOK)
		}
		return

	default:
		// Method not allowed
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
