package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := response{
		Error: false,
		Message: "Hit the broker",
	}

	log.Println("Hit!")

	serialized, _ := json.MarshalIndent(payload, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(serialized)
}
