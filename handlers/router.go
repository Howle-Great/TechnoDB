package handlers

import(
	"github.com/gin-gonic/gin"
)

func CreateRouter() (*gin.Engine) {
	r := gin.Default()

	// Forum part
	forumApi := r.Group("api/forum")
	// forumApi.POST("create", CreateForum)
	forumApi.POST(":slug/create", CreateSlug)
	forumApi.GET(":slug/details", GetForum)
	forumApi.GET(":slug/threads", GetForumThreads)
	forumApi.GET(":slug/users", GetForumUsers)
	forumApi.POST(":slug", CreateForum) // Ready

	// Posts
	postsApi := r.Group("api/post")
	postsApi.GET(":id/details", GetPost)
	postsApi.POST(":id/details", UpdatePost)

	//Services endpoints
	serviceApi := r.Group("api/service")
	serviceApi.POST("clear", Clear)
	serviceApi.GET("status", GetStatus)

	// Thread endpoints
	threadApi := r.Group("api/thread")
	threadApi.POST(":slug_or_id/create", CreatePost)
	threadApi.GET(":slug_or_id/details", GetThread)
	threadApi.POST(":slug_or_id/details", UpdateThread)
	threadApi.GET(":slug_or_id/posts", GetThreadPosts)
	threadApi.POST(":slug_or_id/vote", MakeThreadVote)

	// Users part
	usersApi := r.Group("api/user")
	usersApi.POST(":nickname/create", CreateUser)  
	usersApi.GET(":nickname/profile", GetUser) 
	usersApi.POST(":nickname/profile", UpdateUser) 

	return r
}