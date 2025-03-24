package main

import (
	"log"
	"net/http"
	"os"

	c "github.com/t1001001/prog-assig02/internal/constants"
	h "github.com/t1001001/prog-assig02/internal/handlers"
)

func main() {
	//Setting the port
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Println("Port is not set - Setting port to default: " + port)
	}

	// Setting the endpoints
	http.HandleFunc(c.ROOT+c.VERSION+c.DEFAULT_PATH, h.DefaultHandler)
	http.HandleFunc(c.ROOT+c.VERSION+c.REGISTRATIONS_PATH, h.RegistrationsHandler)
	http.HandleFunc(c.ROOT+c.VERSION+c.DASHBOARDS_PATH, h.DashboardsHandler)
	http.HandleFunc(c.ROOT+c.VERSION+c.NOTIFICATIONS_PATH, h.NotificationsHandler)
	http.HandleFunc(c.ROOT+c.VERSION+c.STATUS_PATH, h.StatusHandler)

	// Starting the server
	log.Println("Starting server on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
