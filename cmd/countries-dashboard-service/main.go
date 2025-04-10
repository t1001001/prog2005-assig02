package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"

	c "github.com/t1001001/prog-assig02/internal/constants"
	h "github.com/t1001001/prog-assig02/internal/handlers"
)

// App struct to hold the Firestore client and other dependencies
type App struct {
	FirestoreClient *firestore.Client
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setting the port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("Port is not set - Using default: " + port)
	}

	// Initialize Firestore
	ctx := context.Background()
	firestoreClient, err := initFirestore(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer firestoreClient.Close()

	// Create an instance of App with Firestore client
	app := &App{FirestoreClient: firestoreClient}

	// Setting the endpoints with Firestore-aware handlers
	http.HandleFunc(c.ROOT+c.VERSION+c.DEFAULT_PATH, app.wrapHandler(h.DefaultHandler))
	http.HandleFunc(c.ROOT+c.VERSION+c.REGISTRATIONS_PATH, app.wrapHandler(h.RegistrationsHandler))
	http.HandleFunc(c.ROOT+c.VERSION+c.DASHBOARDS_PATH, app.wrapHandler(h.DashboardsHandler))
	http.HandleFunc(c.ROOT+c.VERSION+c.NOTIFICATIONS_PATH, app.wrapHandler(h.NotificationsHandler))
	http.HandleFunc(c.ROOT+c.VERSION+c.STATUS_PATH, app.wrapHandler(h.StatusHandler))

	// connecting to the static files
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("internal/assets/images"))))

	// Starting the server
	log.Println("Starting server on port " + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Initialize Firestore with credentials from .env file
func initFirestore(ctx context.Context) (*firestore.Client, error) {
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsPath == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
		return nil, fmt.Errorf("environment variable GOOGLE_APPLICATION_CREDENTIALS is required")
	}

	sa := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	// Test Firestore connection
	_, err = client.Collections(ctx).GetAll()
	if err != nil {
		fmt.Println("Firestore NOT connected properly:", err)
	} else {
		fmt.Println("Firestore connected successfully")
	}

	return client, nil
}

// Middleware to inject Firestore client into handlers
func (app *App) wrapHandler(handlerFunc func(http.ResponseWriter, *http.Request, *firestore.Client)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r, app.FirestoreClient)
	}
}
