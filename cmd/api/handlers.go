package main

import (
	"errors"
	"fmt"
	"net/http"
)


type RequestPayload struct {
	Action            ActionType         `json:"action"`
	UserServiceCreate UserServicePayload `json:"user_service,omitempty"`

	UserServiceCreateViaRabbit UserServiceViaRabbitPayload `json:"user_service_via_rabbit,omitempty"`

	UserServiceGet struct {
		Email string `json:"email"`
	} `json:"user_service_get,omitempty"`

	UserServiceGetAll struct {
	} `json:"user_service_get_all,omitempty"`

	UserServiceGetOne struct {
		ID int `json:"id"`
	} `json:"user_service_get_one,omitempty"`

	PostServiceCreate PostServicePayload `json:"post_service,omitempty"`

	PostServiceCreateViaRabbit PostServiceViaRabbitPayload `json:"post_service_via_rabbit,omitempty"`

	PostServiceGet struct {
		ID int `json:"id"`
	} `json:"post_service_get,omitempty"`

	PostServiceGetAll struct {
	} `json:"post_service_get_all,omitempty"`

	GetPostsRPC GetPostsRPCParams `json:"get_posts_via_rpc,omitempty"`

	GetPostViaGrpc struct {
		ID int `json:"id"`
	} `json:"get_post_via_params_grpc,omitempty"`

	UpdatePostViaRabbit PostUpdateViaRabbitPayload `json:"update_post_via_rabbit,omitempty"`
}




func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case ActionCreateUser:
		app.createUser(w, requestPayload.UserServiceCreate)

	case ActionGetUser:
		app.getUserByID(w, requestPayload.UserServiceGetOne.ID)

	case ActionGetUsers:
		app.getUsers(w)

	case ActionCreatePost:
		app.createPost(w, requestPayload.PostServiceCreate)

	case GetPostsViaRPC:
		app.getAllPostViaRPC(w, requestPayload.GetPostsRPC)

	case ActionCreatePostViaRabbit:
		app.postCreateViaRabbit(w, requestPayload.PostServiceCreateViaRabbit)

	case ActionCreateUserViaRabbit:
		app.userCreateViaRabbit(w, requestPayload.UserServiceCreateViaRabbit)

	case ActionGetPost:
		app.getPostByID(w, requestPayload.PostServiceGet.ID)

	case ActionGetPostViaGrpc:
		app.PostIdViaGrpc(w, r, fmt.Sprintf("%d", requestPayload.GetPostViaGrpc.ID))

	case ActionGetPosts:
		app.getPosts(w)

	case ACTION_UPDATE_POST_VIA_RABBIT:
		app.postUpdateViaRabbit(w, requestPayload.UpdatePostViaRabbit)

	default:
		app.errorJSON(w, errors.New("invalid action"))

	}
}
