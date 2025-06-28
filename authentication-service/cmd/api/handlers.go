package main

import (
	"errors"
	"log"
	"net/http"
)

func (app *Config) Auth(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	//hahahahaviwiovnwvwvwvwnbvioenoibneoib

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		log.Println("Can't read request Payload")
		return
	}

	// validate user
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		log.Println("Can't get user from DB")
		return
	}

	// validate password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("Invalid password"), http.StatusBadRequest)
		log.Println("Invalid password")
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Successfully authenticated",
		Data:    user,
	}

	log.Println(payload)
	log.Println(&payload.Data)

	app.writeJSON(w, http.StatusAccepted, payload)
}
