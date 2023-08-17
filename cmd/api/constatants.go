package main

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
