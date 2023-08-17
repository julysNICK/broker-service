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




func (app *Config) postCreateViaRabbit(w http.ResponseWriter, payload PostServiceViaRabbitPayload) {
	fmt.Println("Creating Post via RabbitMQ")

		payloadConvert := PostServiceViaRabbitPayload{
		Id_user: payload.Id_user,
		Content: payload.Content,
		Type:    payload.Type,
	}


	err := app.pushToQueuePostRabbit(payloadConvert, "post.created")

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

// func (app *Config) pushToQueuePost(content string, id_user int, typeAction string) error {
// 	fmt.Println("Pushing to queue")
// 	emitter, err := event.NewEventEmitterPost(app.RabbitPost)

// 	if err != nil {
// 		fmt.Println("line 409 " + err.Error())
// 		return err
// 	}

// 	fmt.Println("Pushing to channel")

// 	payload := PostServiceViaRabbitPayload{
// 		Id_user: id_user,
// 		Content: content,
// 		Type:    typeAction,
// 	}

// 	j, _ := json.MarshalIndent(payload, "", "\t")
// 	fmt.Println(string(j))

// 	err = emitter.PushPost(string(j), "post.created")

// 	fmt.Println("Pushed to channel")

// 	if err != nil {
// 		fmt.Println("line 428 " + err.Error())
// 		return err
// 	}

// 	return nil
// }

func (app *Config) postUpdateViaRabbit(w http.ResponseWriter, payload PostUpdateViaRabbitPayload) {
	
	fmt.Println("Updating Post via RabbitMQ")

	err := app.pushToQueuePostRabbit(payload, "post.updated")

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
// func (app *Config) pushToQueuePostUpdate( payload PostUpdateViaRabbitPayload) error {
// 	fmt.Println("Pushing to queue")
// 	emitter, err := event.NewEventEmitterPost(app.RabbitPost)

// 	if err != nil {
// 		fmt.Println("line 409 " + err.Error())
// 		return err
// 	}

// 	fmt.Println("Pushing to channel")

	

// 	j, _ := json.MarshalIndent(payload, "", "\t")
// 	fmt.Println(string(j))

// 	err = emitter.PushPost(string(j), "post.updated")

// 	if err != nil {
// 		fmt.Println("line 428 " + err.Error())
// 		return err
// 	}
// 	fmt.Println("Pushed to channel")

// 	return nil
// }


func (app *Config) pushToQueuePostRabbit(payload interface{}, eventType string) error {

	emitter, err := event.NewEventEmitterPost(app.RabbitPost)

	if err != nil {
		fmt.Println("line 264 " + err.Error())
		return err
	}

	j, _ := json.MarshalIndent(payload, "", "\t")

	fmt.Println(string(j))

	err = emitter.PushPost(string(j), eventType)

	if err != nil {
		fmt.Println("line 275 " + err.Error())
		return err
	}

	fmt.Println("Pushed to channel")

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
