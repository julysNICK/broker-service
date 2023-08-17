package main

import (
	"broker/cmd/event"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) createUser(w http.ResponseWriter, payload UserServicePayload) {
	jsonDat, _ := json.MarshalIndent(payload, "", "\t")

	fmt.Println(string(jsonDat))

	request, err := http.NewRequest("POST", "http://user-service/v1/user", bytes.NewBuffer(jsonDat))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		app.errorJSON(w, errors.New("invalid status code"))
		return
	}
	var jsonFromResponse jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromResponse)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromResponse.Error {
		app.errorJSON(w, errors.New(jsonFromResponse.Message))
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Message = "User Created"
	payloadResponse.Error = false
	payloadResponse.Data = jsonFromResponse.Data

	app.writeJSON(w, http.StatusCreated, payloadResponse)

}

func (app *Config) getUserByID(w http.ResponseWriter, id int) {
	request, err := http.NewRequest("GET", "http://user-service/v1/user/"+string(id), nil)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("invalid status code"))
		return
	}
	var jsonFromResponse jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromResponse)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromResponse.Error {
		app.errorJSON(w, errors.New(jsonFromResponse.Message))
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Message = "User Found"
	payloadResponse.Error = false
	payloadResponse.Data = jsonFromResponse.Data

	app.writeJSON(w, http.StatusOK, payloadResponse)
}

func (app *Config) userCreateViaRabbit(w http.ResponseWriter, payload UserServiceViaRabbitPayload) {
	fmt.Println("Creating User via RabbitMQ")
	fmt.Println("line 387 " + payload.Email)
	err := app.pushToQueueUser(UserServiceViaRabbitPayload{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Password:  payload.Password,
		Active:    payload.Active,
		Type:      payload.Type,
	})

	if err != nil {
		fmt.Println("line 390 " + err.Error())
		app.errorJSON(w, err)
		return
	}

	fmt.Println("User Created via RabbitMQ")

	var payloadResponse jsonResponse
	payloadResponse.Message = "User Created"
	payloadResponse.Error = false

	app.writeJSON(w, http.StatusCreated, payloadResponse)
}

func (app *Config) pushToQueueUser(user UserServiceViaRabbitPayload) error {
	fmt.Println("Pushing to queue")
	emitter, err := event.NewEventEmitterUser(app.RabbitPost)

	if err != nil {
		fmt.Println("line 409 " + err.Error())
		return err
	}

	fmt.Println("Pushing to channel")

	payload := UserServiceViaRabbitPayload{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Active:    user.Active,
		Type:      user.Type,
	}

	j, _ := json.MarshalIndent(payload, "", "\t")
	fmt.Println(string(j))

	err = emitter.PushUser(string(j), "user.created")

	fmt.Println("Pushed to channel")

	if err != nil {
		fmt.Println("line 428 " + err.Error())
		return err
	}

	return nil
}

func (app *Config) getUsers(w http.ResponseWriter) {
	request, err := http.NewRequest("GET", "http://user-service/v1/users", nil)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("invalid status code"))
		return
	}
	var jsonFromResponse jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromResponse)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromResponse.Error {
		app.errorJSON(w, errors.New(jsonFromResponse.Message))
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Message = "Users Found"
	payloadResponse.Error = false
	payloadResponse.Data = jsonFromResponse.Data

	app.writeJSON(w, http.StatusOK, payloadResponse)
}