package main

type GetPostsRPCParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
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
	Type    string `json:"type"`
}

type DeletePostViaRabbitPayload struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}
