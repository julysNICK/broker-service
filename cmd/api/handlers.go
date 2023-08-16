package main

import (
	"broker/cmd/event"
	"broker/post"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ActionType string

const (
	ActionCreateUser          ActionType = "create_user"
	ActionCreateUserViaRabbit ActionType = "create_user_via_rabbit"
	GetPostsViaRPC            ActionType = "get_posts_via_rpc"
	ActionGetUser             ActionType = "get_user"
	ActionGetUsers            ActionType = "get_users"

	ActionCreatePost              ActionType = "create_post"
	ActionCreatePostViaRabbit     ActionType = "create_post_via_rabbit"
	ACTION_UPDATE_POST_VIA_RABBIT ActionType = "update_post_via_rabbit"
	ActionGetPost                 ActionType = "get_post"
	ActionGetPosts                ActionType = "get_posts"
	ActionGetPostViaGrpc          ActionType = "get_post_via_grpc"
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

type UserServicePayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Active    int    `json:"active"`
}

type GetPostsRPCParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type UserServiceViaRabbitPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Active    int    `json:"active"`
	Type      string `json:"type"`
}

type PostServicePayload struct {
	Id_user int    `json:"id_user"`
	Content string `json:"content"`
}

type PostServiceViaRabbitPayload struct {
	Id_user int    `json:"id_user"`
	Content string `json:"content"`
	Type    string `json:"type"`
}


type PostUpdateServiceViaRabbitPayload struct {
	Id_user int    `json:"id_user"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type PostUpdateViaRabbitPayload struct {
	Id_user int    `json:"id_user"`
	Content string `json:"content"`
	Type 	string `json:"type"`
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

func (app *Config) getPosts(w http.ResponseWriter) {
	request, err := http.NewRequest("GET", "http://post-service/v1/posts", nil)

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
	payloadResponse.Message = "Posts Found"
	payloadResponse.Error = false
	payloadResponse.Data = jsonFromResponse.Data

	app.writeJSON(w, http.StatusOK, payloadResponse)
}

func (app *Config) getPostByID(w http.ResponseWriter, id int) {
	request, err := http.NewRequest("GET", "http://post-service/v1/post/"+string(id), nil)

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
	payloadResponse.Message = "Post Found"
	payloadResponse.Error = false
	payloadResponse.Data = jsonFromResponse.Data

	app.writeJSON(w, http.StatusOK, payloadResponse)
}

func (app *Config) createPost(w http.ResponseWriter, payload PostServicePayload) {
	jsonDat, _ := json.MarshalIndent(payload, "", "\t")

	request, err := http.NewRequest("POST", "http://post-service/v1/post", bytes.NewBuffer(jsonDat))

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
	payloadResponse.Message = "Post Created"
	payloadResponse.Error = false
	payloadResponse.Data = jsonFromResponse.Data

	app.writeJSON(w, http.StatusCreated, payloadResponse)

}

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

func (app *Config) postCreateViaRabbit(w http.ResponseWriter, payload PostServiceViaRabbitPayload) {
	fmt.Println("Creating Post via RabbitMQ")

	err := app.pushToQueuePost(payload.Content, payload.Id_user, payload.Type)

	if err != nil {
		fmt.Println("line 390 " + err.Error())
		app.errorJSON(w, err)
		return
	}

	fmt.Println("Post Created via RabbitMQ")

	var payloadResponse jsonResponse
	payloadResponse.Message = "Post Created"
	payloadResponse.Error = false

	app.writeJSON(w, http.StatusCreated, payloadResponse)
}

func (app *Config) pushToQueuePost(content string, id_user int, typeAction string) error {
	fmt.Println("Pushing to queue")
	emitter, err := event.NewEventEmitterPost(app.RabbitPost)

	if err != nil {
		fmt.Println("line 409 " + err.Error())
		return err
	}

	fmt.Println("Pushing to channel")

	payload := PostServiceViaRabbitPayload{
		Id_user: id_user,
		Content: content,
		Type:    typeAction,
	}

	j, _ := json.MarshalIndent(payload, "", "\t")
	fmt.Println(string(j))

	err = emitter.PushPost(string(j), "post.created")

	fmt.Println("Pushed to channel")

	if err != nil {
		fmt.Println("line 428 " + err.Error())
		return err
	}

	return nil
}


func (app *Config) pushToQueuePostUpdate( payload PostUpdateViaRabbitPayload) error {
	fmt.Println("Pushing to queue")
	emitter, err := event.NewEventEmitterPost(app.RabbitPost)

	if err != nil {
		fmt.Println("line 409 " + err.Error())
		return err
	}

	fmt.Println("Pushing to channel")

	payloadUpdate := PostUpdateServiceViaRabbitPayload{
		Id_user: payload.Id_user,
		Content: payload.Content,
		Type:    payload.Content,
	}

	j, _ := json.MarshalIndent(payloadUpdate, "", "\t")
	fmt.Println(string(j))

	err = emitter.PushPost(string(j), "post.update")

	fmt.Println("Pushed to channel")

	if err != nil {
		fmt.Println("line 428 " + err.Error())
		return err
	}

	return nil
}


func (app *Config) postUpdateViaRabbit(w http.ResponseWriter, payload PostUpdateViaRabbitPayload) {
	
	fmt.Println("Updating Post via RabbitMQ")

	err := app.pushToQueuePostUpdate(payload)

	if err != nil {
		fmt.Println("line 549 " + err.Error())
		app.errorJSON(w, err)
		return
	}

	fmt.Println("Post Updated via RabbitMQ")

	var payloadResponse jsonResponse
	payloadResponse.Message = "Post Updated"
	payloadResponse.Error = false

	app.writeJSON(w, http.StatusCreated, payloadResponse)
}

// -----------------------------------------------------------------------------------------------------------------------------------------------------------------------------
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

type RPCPost struct {
	ID      int
	Title   string
	Content string
}

func (app *Config) getAllPostViaRPC(w http.ResponseWriter, payload GetPostsRPCParams) {
	fmt.Println("Getting Posts via RPC")
	client, err := rpc.Dial("tcp", "post-service:5001")

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var reply []RPCPost

	fmt.Println("Getting Posts via RPC 5001")

	fmt.Println(payload)

	err = client.Call("API.GetPostsRPC", payload, &reply)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Message = "Posts Found"
	payloadResponse.Error = false
	payloadResponse.Data = reply

	app.writeJSON(w, http.StatusOK, payloadResponse)
}

func (app *Config) PostIdViaGrpc(w http.ResponseWriter, r *http.Request, id string) {

	conn, err := grpc.Dial("post-service:5002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	p := post.NewPostServiceClient(conn)

	req := &post.PostRequest{
		Id: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	res, err := p.GetPostRpc(ctx, req)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payloadResponse jsonResponse
	payloadResponse.Message = "Post Found"
	payloadResponse.Error = false
	payloadResponse.Data = res

	app.writeJSON(w, http.StatusOK, payloadResponse)
}
