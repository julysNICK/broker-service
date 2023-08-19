package main

import (
	"broker/cmd/event"
	userGrpc "broker/user"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (app *Config) getAllUsersViaGRPC(w http.ResponseWriter) {

	conn, err := grpc.Dial("user-service:5003", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	u := userGrpc.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := u.GetAllUsers(ctx, &emptypb.Empty{})

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Message = "Users Found"
	payloadResponse.Error = false
	payloadResponse.Data = res

	app.writeJSON(w, http.StatusOK, payloadResponse)
}

func (app *Config) getUserViaGRPC(w http.ResponseWriter, r *http.Request, id string) {
	conn, err := grpc.Dial("user-service:5003", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()
	u := userGrpc.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	convID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	res, err := u.GetOneUser(ctx, &userGrpc.UserRequestGetOne{
		Id: int32(convID),
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	var payloadResponse jsonResponse
	payloadResponse.Message = "User Found"
	payloadResponse.Error = false
	payloadResponse.Data = res
	app.writeJSON(w, http.StatusOK, payloadResponse)
}

func (app *Config) userDeleteViaGRPC(w http.ResponseWriter, email string) {
	conn, err := grpc.Dial("user-service:5003", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	u := userGrpc.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := u.DeleteUser(ctx, &userGrpc.UserRequestDelete{
		Email: email,
	})

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Message = "User Deleted"
	payloadResponse.Error = false
	payloadResponse.Data = res

	app.writeJSON(w, http.StatusOK, payloadResponse)
}

type UserUpdateViaGRPCPayload struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (app *Config) userUpdateViaGRPC(w http.ResponseWriter, payload UserUpdateViaGRPCPayload) {
	conn, err := grpc.Dial("user-service:5003", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	u := userGrpc.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := u.UpdateUser(ctx, &userGrpc.UserRequestUpdate{
		Email: payload.Email,
	})

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payloadResponse jsonResponse

	payloadResponse.Message = "User Updated"

	payloadResponse.Error = false

	payloadResponse.Data = res

	app.writeJSON(w, http.StatusOK, payloadResponse)
}
