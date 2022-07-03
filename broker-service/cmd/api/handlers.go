package main

import (
	"log"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := response{
		Error: false,
		Message: "Hit the broker",
	}

	log.Println("Hit!")

	err := app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Println("broker write JSON error: ", err)
	}
}
