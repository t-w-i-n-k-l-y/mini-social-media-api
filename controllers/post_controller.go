package controllers

import (
	"mini-social-media-api/models"
	"mini-social-media-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CreatePostHandler handles the creation of a new post
// Expects a JSON payload with `content` in the request body
// Returns the created post or an error if the request is invalid or creation fails
func CreatePostHandler(c *gin.Context) {
	var req models.Post
	err := c.ShouldBindJSON(&req) // Parse the request body into the Post model
	if err != nil {
		logrus.Errorln("Failed to create the post: Invalid request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Content should be within 1-250 characters"})
		return
	}

	// Call the service to create a new post
	post, err := services.CreatePost(req.Content)
	if err != nil {
		logrus.Errorln("Failed to create the post: Error occurred in create post service: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to create post: " + err.Error()})
		return
	}

	logrus.Infoln("Post created successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post": post})
}

// UpdatePostHandler handles updating an existing post
// Expects a `postID` as a URL parameter and `content` in the JSON payload
// Returns the updated post or an error if the post is not found or the request is invalid
func UpdatePostHandler(c *gin.Context) {
	postIDParam := c.Param("postID")         // Retrieve the post ID from URL parameters
	postID, err := strconv.Atoi(postIDParam) // Convert post ID to integer
	if err != nil {
		logrus.Errorln("Failed to update post: Error in converting post ID to int")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req models.Post

	err = c.ShouldBindJSON(&req) // Parse the request body into the Post model
	if err != nil {
		logrus.Errorln("Failed to update post: Error in request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Content should be within 1-250 characters"})
		return
	}

	// Call the service to update the post
	post, err := services.UpdatePost(postID, req.Content)
	if err != nil {
		logrus.Errorln("Failed to update post: Error occurred in update post service: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to update post: " + err.Error()})
		return
	}

	logrus.Infoln("Post updated success.  ID: " + postIDParam)
	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully", "post": post})
}

// LikePostHandler increments the like count for a specific post
// Expects a `postID` as a URL parameter
// Returns the updated post or an error if the post is not found or the ID is invalid
func LikePostHandler(c *gin.Context) {
	postIDParam := c.Param("postID")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		logrus.Errorln("Failed to like the post: Error in converting post ID to int: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Call the service to increment likes
	post, err := services.LikePost(postID)
	if err != nil {
		logrus.Errorln("Failed to like the post: Error occurred in like post service: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to like: " + err.Error()})
		return
	}

	logrus.Infoln("Like added successfully. ID: " + postIDParam)
	c.JSON(http.StatusOK, gin.H{"message": "Liked the post successfully", "post": post})
}

// GetPostDetailsHandler retrieves the details of a specific post by ID
// Expects a `postID` as a URL parameter
// Returns the post details or an error if the post is not found
func GetPostDetailsHandler(c *gin.Context) {
	postIDParam := c.Param("postID")
	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		logrus.Errorln("Failed to get details of the post: Error in converting post id to int: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get post details. Invalid post ID"})
		return
	}

	// Fetch post details
	post, err := services.GetPostDetailsByID(postID)
	if err != nil {
		logrus.Errorln("Failed to get post details: Error occurred in get post details service: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to get post details: " + err.Error()})
		return
	}

	// Construct and return the response
	logrus.Infoln("Retrieved post successfully. ID: " + postIDParam)
	c.JSON(http.StatusOK, gin.H{"post": post})
}

// AddCommentHandler adds a comment to a specific post
// Expects a `postID` as a URL parameter and comment text in the JSON payload
// Returns the updated post or an error if the post is not found or the request is invalid
func AddCommentHandler(c *gin.Context) {
	postID := c.Param("postID")
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		logrus.Errorln("Failed to add comment: Error in converting post id to int: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add comment: No post found"})
		return
	}

	var reqComment models.Comment
	err = c.ShouldBindJSON(&reqComment)
	if err != nil {
		logrus.Errorln("Failed to add comment: Error in request body: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. Text should be within 1-150 characters"})
		return
	}

	// Call the service to add a comment
	updatedPost, err := services.AddComment(postIDInt, reqComment)
	if err != nil {
		logrus.Errorln("Failed to add comment: Error occurred in add comment service")
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to add the comment: " + err.Error()})
		return
	}

	logrus.Infof("Comment added to post %v successfully \n", postID)
	c.JSON(http.StatusOK, gin.H{"message": "Comment added successfully", "post": updatedPost})
}

// GetAllPostsHandlerWithPagination retrieves all posts from the in-memory storage with pagination support
// Returns the paginated list of posts
func GetAllPostsHandlerWithPagination(c *gin.Context) {
	// Default pagination parameters if not set
	page := 1
	limit := 10

	// Parse query parameters for pagination
	pg, exists := c.GetQuery("page")
	if exists {
		parsedPage, err := strconv.Atoi(pg)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		} else {
			logrus.Warnln("Invalid page query parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}
	}

	if lt, exists := c.GetQuery("limit"); exists {
		parsedLimit, err := strconv.Atoi(lt)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		} else {
			logrus.Warnln("Invalid limit query parameter")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}

	// Get all posts from the service
	posts := services.GetAllPosts()
	totalPosts := len(posts)

	if totalPosts == 0 {
		logrus.Infoln("Failed to retrieve posts: No posts found")
		c.JSON(http.StatusOK, gin.H{"message": "No posts found"})
		return
	}

	// Calculate pagination boundaries
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex >= totalPosts {
		logrus.Infoln("Page out of range: No posts found")
		c.JSON(http.StatusOK, gin.H{
			"message": "No posts found",
			"page":    page,
			"limit":   limit,
			"total":   totalPosts,
		})
		return
	}

	if endIndex > totalPosts {
		endIndex = totalPosts
	}

	// Paginate posts
	paginatedPosts := posts[startIndex:endIndex]

	logrus.Infof("Retrieved posts for page %d with limit %d", page, limit)
	c.JSON(http.StatusOK, gin.H{
		"posts": paginatedPosts,
		"page":  page,
		"limit": limit,
		"total": totalPosts,
	})
}
