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
	case http.MethodGet:
		handleGetWebhook(w, r, client)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleRegisterWebhook(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	var webhook c.Webhook
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	webhook.ID = uuid.New().String()

	_, err := client.Collection("webhooks").Doc(webhook.ID).Set(context.Background(), webhook)
	if err != nil {
		http.Error(w, "Failed to store webhook", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": webhook.ID})
}

func handleDeleteWebhook(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	id := extractWebhookID(r)
	if id == "" {
		http.Error(w, "Missing webhook ID in URL path", http.StatusBadRequest)
		return
	}

	_, err := client.Collection("webhooks").Doc(id).Delete(context.Background())
	if err != nil {
		http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleGetWebhook(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	id := extractWebhookID(r)

	// Return one webhook if ID is present
	if id != "" {
		doc, err := client.Collection("webhooks").Doc(id).Get(context.Background())
		if err != nil {
			http.Error(w, "Webhook not found", http.StatusNotFound)
			return
		}

		var webhook c.Webhook
		if err := doc.DataTo(&webhook); err != nil {
			http.Error(w, "Failed to parse webhook data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(webhook)
		return
	}

	// Otherwise return all webhooks
	iter := client.Collection("webhooks").Documents(context.Background())
	var webhooks []c.Webhook
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var wh c.Webhook
		if err := doc.DataTo(&wh); err == nil {
			webhooks = append(webhooks, wh)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhooks)
}

func extractWebhookID(r *http.Request) string {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 5 {
		return parts[4]
	}
	return ""
}
