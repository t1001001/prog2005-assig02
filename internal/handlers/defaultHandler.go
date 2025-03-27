package handlers

import (
	"net/http"

	"cloud.google.com/go/firestore"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request, client *firestore.Client) {

}
