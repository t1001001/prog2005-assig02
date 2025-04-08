package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/t1001001/prog-assig02/internal/constants"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

// NotificationsHandler handles POST requests to register new webhooks
func NotificationsHandler(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var webhook constants.Webhook
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Generate unique ID
	webhook.ID = uuid.New().String()

	// Store in Firestore
	_, err = client.Collection("webhooks").Doc(webhook.ID).Set(context.Background(), webhook)
	if err != nil {
		http.Error(w, "Failed to store webhook", http.StatusInternalServerError)
		return
	}

	// Respond with ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": webhook.ID})
}
