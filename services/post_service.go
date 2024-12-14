package services

import (
	"errors"
	"mini-social-media-api/models"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var posts []models.Post       // In-memory storage for all posts
var postMutex = &sync.Mutex{} // Mutex to ensure safe concurrent access to the posts slice
var postIDCounter = 1         // Counter for generating unique post IDs
var now = time.Now().Local()  // Current local time

// CreatePost creates a new post with the given content.
// Returns the created post or an error if the content is invalid.
func CreatePost(content string) (models.Post, error) {
	// Validate content
	if content == "" || strings.TrimSpace(content) == "" {
		return models.Post{}, errors.New("post content cannot be empty")
	}
	if len(content) > 250 {
		return models.Post{}, errors.New("post content exceeds maximum length of 250 characters")
	}

	postMutex.Lock()
	defer postMutex.Unlock()

	// Initialize a new post with default values and given content
	post := models.Post{
		ID:        postIDCounter,
		Content:   content,
		Likes:     0,
		Comments:  []models.Comment{},
		CreatedAt: now,
		UpdatedAt: now,
	}

	postIDCounter++             // Increment the counter for the next postID
	posts = append(posts, post) // Add the new post to the in-memory slice

	return post, nil
}

// UpdatePost updates the content of an existing post by its ID.
// Returns the updated post or an error if the post is not found or the content is invalid.
func UpdatePost(id int, newContent string) (models.Post, error) {
	postMutex.Lock()
	defer postMutex.Unlock()

	// Validate the new content
	if newContent == "" || strings.TrimSpace(newContent) == "" {
		return models.Post{}, errors.New("post content cannot be empty")
	}
	if len(newContent) > 250 {
		return models.Post{}, errors.New("post content exceeds maximum length of 250 characters")
	}

	// Find the post by ID and update its content
	for i, post := range posts {
		if post.ID == id {
			posts[i].Content = newContent
			posts[i].UpdatedAt = time.Now()
			return posts[i], nil
		}
	}

	return models.Post{}, errors.New("post not found")
}

// GetAllPosts retrieves all posts from the in-memory storage.
// Returns a slice of all posts.
func GetAllPosts() []models.Post {
	postMutex.Lock()
	defer postMutex.Unlock()
	return posts
}

// LikePost increments the like count for a specific post by its ID.
// Returns the updated post or an error if the post is not found.
func LikePost(id int) (models.Post, error) {
	postMutex.Lock()
	defer postMutex.Unlock()

	// Find the post by its ID and increment its like count
	for i, post := range posts {
		if post.ID == id {
			posts[i].Likes++
			return posts[i], nil
		}
	}

	return models.Post{}, errors.New("post not found")
}

// GetPostDetailsByID retrieves the details of a specific post by its ID, excluding comments.
// Returns the found post or an error if the post is not found.
func GetPostDetailsByID(id int) (models.Post, error) {
	postMutex.Lock()
	defer postMutex.Unlock()

	// Iterate over posts to find the one matching the given ID
	for _, post := range posts {
		if post.ID == id {
			return post, nil
		}
	}

	return models.Post{}, errors.New("post not found")
}

// AddComment adds a new comment to a specific post by its ID.
// Returns the updated post or an error if the post is not found or validation fails.
func AddComment(postID int, comment models.Comment) (models.Post, error) {
	// Validate comment text
	if comment.Text == "" || strings.TrimSpace(comment.Text) == "" {
		logrus.Error()
		return models.Post{}, errors.New("comment cannot be empty")
	}
	if len(comment.Text) > 150 {
		return models.Post{}, errors.New("comment exceeds maximum length of 150 characters")
	}

	postMutex.Lock()
	defer postMutex.Unlock()

	// find the post by ID
	for i, post := range posts {
		if post.ID == postID {
			// Create a new comment
			newComment := models.Comment{
				ID:        len(post.Comments) + 1, // Generate comment ID based on the length of the Comments slice
				Text:      comment.Text,
				CreatedAt: now,
			}

			// Append the new comment to the post's comments slice
			posts[i].Comments = append(posts[i].Comments, newComment)

			return posts[i], nil
		}
	}

	return models.Post{}, errors.New("post not found")
}
