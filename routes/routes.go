package routes

import (
	"mini-social-media-api/controllers"

	"github.com/gin-gonic/gin"
)

// InitRoutes initializes all the application routes and returns the configured Gin router
func InitRoutes() *gin.Engine {
	router := gin.Default()

	// Grouping routes related to posts for better organization
	postRoutes := router.Group("/posts")
	{
		postRoutes.POST("/", controllers.CreatePostHandler)                 // Route to create a new post
		postRoutes.PUT("/:postID", controllers.UpdatePostHandler)           // Route to update an existing post
		postRoutes.GET("/", controllers.GetAllPostsHandlerWithPagination)   // Route to get all posts
		postRoutes.GET("/:postID", controllers.GetPostDetailsHandler)       // Route to get details of a specific post by ID
		postRoutes.POST("/:postID/like", controllers.LikePostHandler)       // Route to like a specific post
		postRoutes.POST("/:postID/comments", controllers.AddCommentHandler) // Route to add a comment to a specific post
	}

	return router
}
