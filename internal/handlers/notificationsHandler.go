package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"

	c "github.com/t1001001/prog-assig02/internal/constants"
)

// NotificationsHandler handles webhook requests
func NotificationsHandler(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	switch r.Method {
	case http.MethodPost:
		handleRegisterWebhook(w, r, client)
	case http.MethodDelete:
		handleDeleteWebhook(w, r, client)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleRegisterWebhook(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	var webhook c.Webhook
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Generate a unique ID
	webhook.ID = uuid.New().String()

	// Store the webhook in Firestore
	_, err = client.Collection("webhooks").Doc(webhook.ID).Set(context.Background(), webhook)
	if err != nil {
		http.Error(w, "Failed to store webhook", http.StatusInternalServerError)
		return
	}

	// Respond with the webhook ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": webhook.ID})
}

func handleDeleteWebhook(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	// Extract the ID from the path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 || parts[4] == "" {
		http.Error(w, "Missing webhook ID in URL path", http.StatusBadRequest)
		return
	}
	id := parts[4]

	// Delete the webhook from Firestore
	_, err := client.Collection("webhooks").Doc(id).Delete(context.Background())
	if err != nil {
		http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
		return
	}

	// Success
	w.WriteHeader(http.StatusNoContent)
}
