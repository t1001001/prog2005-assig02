package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"

	c "github.com/t1001001/prog-assig02/internal/constants"
	h "github.com/t1001001/prog-assig02/internal/helpers"
)

// respondWithJSON sends a pretty-formatted JSON response
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

// RegistrationsHandler handles all HTTP methods for the /registrations endpoint
func RegistrationsHandler(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	ctx := context.Background()
	path := strings.TrimPrefix(r.URL.Path, c.ROOT+c.VERSION+c.REGISTRATIONS_PATH)
	id := strings.Trim(path, "/")

	switch r.Method {
	case http.MethodPost:
		handlePostRegistration(w, r, ctx, client)
	case http.MethodGet:
		handleGetRegistration(w, r, ctx, client, id)
	case http.MethodPut:
		handlePutRegistration(w, r, ctx, client, id)
	case http.MethodDelete:
		handleDeleteRegistration(w, r, ctx, client, id)
	case http.MethodHead:
		handleHeadRegistration(w, r, ctx, client, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// POST: Register new dashboard configuration
func handlePostRegistration(w http.ResponseWriter, r *http.Request, ctx context.Context, client *firestore.Client) {
	var config c.Country
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if config.Country == "" && config.ISOCode == "" {
		http.Error(w, "Either 'country' or 'isoCode' must be provided", http.StatusBadRequest)
		return
	}

	config.LastRetrieval = h.GetCurrentTimestamp()
	id := uuid.New().String()

	if _, err := client.Collection("registrations").Doc(id).Set(ctx, config); err != nil {
		http.Error(w, "Failed to store configuration", http.StatusInternalServerError)
		return
	}

	log.Printf("[POST] Registered config for: %s (%s)", config.Country, config.ISOCode)
	respondWithJSON(w, http.StatusOK, map[string]string{
		"id":         id,
		"lastChange": config.LastRetrieval,
	})
}

// GET: Retrieve one or all registrations
func handleGetRegistration(w http.ResponseWriter, r *http.Request, ctx context.Context, client *firestore.Client, id string) {
	if id != "" {
		doc, err := client.Collection("registrations").Doc(id).Get(ctx)
		if err != nil {
			http.Error(w, "Configuration not found", http.StatusNotFound)
			return
		}
		var config c.Country
		if err := doc.DataTo(&config); err != nil {
			http.Error(w, "Error parsing configuration", http.StatusInternalServerError)
			return
		}
		response := c.RegistrationResponse{
			ID:         doc.Ref.ID,
			Country:    config.Country,
			ISOCode:    config.ISOCode,
			Features:   config.Features,
			LastChange: config.LastRetrieval,
		}
		respondWithJSON(w, http.StatusOK, response)
	} else {
		docs, err := client.Collection("registrations").Documents(ctx).GetAll()
		if err != nil {
			http.Error(w, "Failed to fetch configurations", http.StatusInternalServerError)
			return
		}
		var configs []c.RegistrationResponse
		for _, doc := range docs {
			var config c.Country
			doc.DataTo(&config)
			entry := c.RegistrationResponse{
				ID:         doc.Ref.ID,
				Country:    config.Country,
				ISOCode:    config.ISOCode,
				Features:   config.Features,
				LastChange: config.LastRetrieval,
			}
			configs = append(configs, entry)
		}
		respondWithJSON(w, http.StatusOK, configs)
	}
}

// PUT: Update a specific configuration
func handlePutRegistration(w http.ResponseWriter, r *http.Request, ctx context.Context, client *firestore.Client, id string) {
	if id == "" {
		http.Error(w, "ID required in URL", http.StatusBadRequest)
		return
	}

	var config c.Country
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config.LastRetrieval = h.GetCurrentTimestamp()

	if _, err := client.Collection("registrations").Doc(id).Set(ctx, config); err != nil {
		http.Error(w, "Failed to update configuration", http.StatusInternalServerError)
		return
	}

	log.Printf("[PUT] Updated config for ID: %s", id)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE: Remove a configuration
func handleDeleteRegistration(w http.ResponseWriter, r *http.Request, ctx context.Context, client *firestore.Client, id string) {
	if id == "" {
		http.Error(w, "ID required in URL", http.StatusBadRequest)
		return
	}

	if _, err := client.Collection("registrations").Doc(id).Delete(ctx); err != nil {
		http.Error(w, "Failed to delete configuration", http.StatusInternalServerError)
		return
	}

	log.Printf("[DELETE] Deleted config for ID: %s", id)
	w.WriteHeader(http.StatusNoContent)
}

// HEAD: Check if config or collection exists
func handleHeadRegistration(w http.ResponseWriter, r *http.Request, ctx context.Context, client *firestore.Client, id string) {
	if id != "" {
		_, err := client.Collection("registrations").Doc(id).Get(ctx)
		if err != nil {
			http.Error(w, "Configuration not found", http.StatusNotFound)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
