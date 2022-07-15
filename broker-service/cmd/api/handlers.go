package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := response{
		Error:   false,
		Message: "Hit the broker",
	}

	log.Println("Hit!")

	err := app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Println("broker write JSON error: ", err)
	}
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	authResponse, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer authResponse.Body.Close()

	if authResponse.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if authResponse.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	var authReponseJson response
	err = json.NewDecoder(authResponse.Body).Decode(&authReponseJson)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if authReponseJson.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload response
	payload.Error = false
	payload.Message = "authenticated"
	payload.Data = authReponseJson.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}
