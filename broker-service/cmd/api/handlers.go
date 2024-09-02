package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type requestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var reqPayload requestPayload

	err := app.readJSON(w, r, &reqPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch reqPayload.Action {
	case "auth":
		app.Authenticate(w, reqPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}

}

func (app *Config) Authenticate(w http.ResponseWriter, a AuthPayload) {
	// Create JSON will send to auth microservice
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		log.Println("Can't convert to JSON")
		return
	}

	// Call service
	request, err := http.NewRequest("POST", "http://authentication-service:8082/auth", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		log.Println(response.Body)
		return
	}
	defer response.Body.Close()
	log.Println("Checking status codes")

	// Make sure to get valid status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling AUTH SERVICE"))
		log.Println(response.StatusCode)
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	log.Println(jsonFromService.Data)

	if err != nil {
		app.errorJSON(w, err)
		log.Println(jsonFromService)
		log.Println("Error in Decode")
		return
	}

	if jsonFromService.Error {
		log.Println("Error in jsonFromService")
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}
