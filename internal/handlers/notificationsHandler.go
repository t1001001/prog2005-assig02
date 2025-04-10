package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

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

	if webhook.URL == "" || webhook.Event == "" {
		http.Error(w, "Missing required fields: url or event", http.StatusBadRequest)
		return
	}

	// Validate event type
	switch webhook.Event {
	case c.EVENT_REGISTER, c.EVENT_CHANGE, c.EVENT_DELETE, c.EVENT_INVOKE:
	default:
		http.Error(w, "Invalid event type", http.StatusBadRequest)
		return
	}

	// Set the current time when the webhook is registered
	webhook.ID = uuid.New().String()
	webhook.Time = time.Now().Format("20060102 15:04")

	_, err := client.Collection("webhooks").Doc(webhook.ID).Set(context.Background(), webhook)
	if err != nil {
		http.Error(w, "Failed to store webhook", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": webhook.ID, "time": webhook.Time})
}

func handleDeleteWebhook(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	id := extractWebhookID(r)
	if id == "" {
		http.Error(w, "Missing webhook ID in URL path", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	docRef := client.Collection("webhooks").Doc(id)

	// Check if document exists
	_, err := docRef.Get(ctx)
	if err != nil {
		http.Error(w, "Webhook not found", http.StatusNotFound)
		return
	}

	// Delete the document
	_, err = docRef.Delete(ctx)
	if err != nil {
		http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Webhook deleted successfully",
		"id":      id,
	})
}

func handleGetWebhook(w http.ResponseWriter, r *http.Request, client *firestore.Client) {
	id := extractWebhookID(r)

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
	if len(parts) >= 4 && parts[2] == "notifications" {
		return parts[3]
	}
	return ""
}

// sendWebhookNotification sends an HTTP POST to a registered webhook URL
func sendWebhookNotification(webhook c.Webhook, countryCode string, event string) {
	// Log the details of the webhook invocation attempt
	log.Printf("Invoking webhook for event '%s' to URL: %s with countryCode: %s", event, webhook.URL, countryCode)

	// Prepare the payload for the webhook
	payload := map[string]string{
		"id":      webhook.ID,
		"country": countryCode,
		"event":   event,
		"time":    time.Now().Format("20060102 15:04"), // Format: YYYYMMDD HH:mm
	}

	// Marshal the payload to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal webhook payload for ID %s: %v", webhook.ID, err)
		return
	}

	// Send the HTTP POST request to the webhook URL
	resp, err := http.Post(webhook.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to send webhook ID %s to %s: %v", webhook.ID, webhook.URL, err)
		return
	}
	defer resp.Body.Close()

	// Log the status code of the response
	if resp.StatusCode == http.StatusOK {
		log.Printf("Webhook %s triggered to %s with event '%s' - Status: %s", webhook.ID, webhook.URL, event, resp.Status)
	} else {
		log.Printf("Webhook %s triggered to %s with event '%s' - Status: %s, Body: %s", webhook.ID, webhook.URL, event, resp.Status, resp.Status)
	}
}

// TriggerWebhooks sends notifications to all matching webhooks
func TriggerWebhooks(client *firestore.Client, countryCode, event string) {
	ctx := context.Background()
	iter := client.Collection("webhooks").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		var webhook c.Webhook
		if err := doc.DataTo(&webhook); err == nil {
			// Match event type and optionally country
			if webhook.Event == event && (webhook.Country == "" || webhook.Country == countryCode) {
				go sendWebhookNotification(webhook, countryCode, event) // This will log when invoked
			}
		}
	}
}
