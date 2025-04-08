package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/t1001001/prog-assig02/internal/constants"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

// NotificationsHandler handles webhook registration
func NotificationsHandler(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	switch r.Method {
	case http.MethodPost:
		var webhook constants.Webhook

		// Decode request body
		err := json.NewDecoder(r.Body).Decode(&webhook)
		if err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if webhook.URL == "" || webhook.Event == "" {
			http.Error(w, "Missing required fields: url and event", http.StatusBadRequest)
			return
		}

		// Generate a unique ID for this webhook
		webhook.ID = uuid.New().String()

		// Save to Firestore
		_, err = client.Collection("webhooks").Doc(webhook.ID).Set(context.Background(), webhook)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save webhook: %v", err), http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(webhook)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
