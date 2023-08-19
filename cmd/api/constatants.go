package main

type ActionType string

const (
	ActionCreateUser          ActionType = "create_user"
	ActionGetUser             ActionType = "get_user"
	ActionGetUsers            ActionType = "get_users"
	ActionCreateUserViaRabbit ActionType = "create_user_via_rabbit"
	ActionGetUsersViaGRPC     ActionType = "get_users_via_grpc"
	ActionGetUserViaGRPC      ActionType = "get_user_via_grpc"
	ActionDeleteUserViaGRPC   ActionType = "delete_user_via_grpc"
	ActionUpdateViaGRPC       ActionType = "update_user_via_grpc"

	GetPostsViaRPC                ActionType = "get_posts_via_rpc"
	ActionCreatePost              ActionType = "create_post"
	ActionCreatePostViaRabbit     ActionType = "create_post_via_rabbit"
	ACTION_UPDATE_POST_VIA_RABBIT ActionType = "update_post_via_rabbit"
	ActionGetPost                 ActionType = "get_post"
	ActionGetPosts                ActionType = "get_posts"
	ActionGetPostViaGrpc          ActionType = "get_post_via_grpc"
	ActionDeletePostViaRabbit     ActionType = "delete_post_via_rabbit"
)
