package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	d, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "error serialization JSON", http.StatusInternalServerError)
	}

	_, err = w.Write(d)
	if err != nil {
		log.Printf("error writing response: %w", err)
		return
	}
}
