package services

import (
	"mini-social-media-api/models"
	"strings"
	"testing"
	"time"
)

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantErr    bool
		wantLikes  int
		wantLength int
	}{
		// Valid cases
		{"Valid content", "This is a valid post", false, 0, 20},                       // Expected: No error, Likes=0, Length=20
		{"Min length content", "a", false, 0, 1},                                      // Expected: No error, Likes=0, Length=1 (Min Length)
		{"Max length content", strings.Repeat("a", 250), false, 0, 250},               // Expected: No error, Likes=0, Length=250 (Max Length)
		{"Content with leading/trailing spaces", "   valid content   ", false, 0, 19}, // Expected: No error, Likes=0, Length=19

		// Invalid cases
		{"Empty content", "", true, 0, 0},                          // Expected: Error, no valid post
		{"Too long content", strings.Repeat("a", 251), true, 0, 0}, // Expected: Error, no valid post
		{"Content with white spaces", "     ", true, 0, 0},         // Expected: Error, no valid post
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			post, err := CreatePost(testCase.content)

			// Check error matches expected result
			if (err != nil) != testCase.wantErr {
				t.Fatalf("Test '%s' failed: expected error = %v, got = %v", testCase.name, testCase.wantErr, err)
			}

			// If no error is expected, check the post's properties
			if !testCase.wantErr {
				// Validate content length
				if len(post.Content) != testCase.wantLength {
					t.Errorf("Test '%s' failed: expected content length = %d, got = %d", testCase.name, testCase.wantLength, len(post.Content))
				}

				// Validate likes
				if post.Likes != testCase.wantLikes {
					t.Errorf("Test '%s' failed: expected likes = %d, got = %d", testCase.name, testCase.wantLikes, post.Likes)
				}

				// Validate content matches
				if post.Content != testCase.content {
					t.Errorf("Test '%s' failed: expected content = '%s', got = '%s'", testCase.name, testCase.content, post.Content)
				}
			} else {
				// If an error is expected, ensure post is empty
				if post.Content != "" || post.Likes != 0 {
					t.Errorf("Test '%s' failed: post should be empty, got: %+v", testCase.name, post)
				}
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	// Setup initial data
	initialPost := models.Post{
		ID:        1,
		Content:   "Initial Content",
		Likes:     0,
		Comments:  []models.Comment{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	posts = []models.Post{initialPost}

	tests := []struct {
		name       string
		id         int
		newContent string
		wantErr    bool
	}{
		// Valid cases
		{"Valid update", 1, "Updated Content", false},
		{"Update with min length content", 1, "a", false},                          // Edge case: minimum length
		{"Update with max length content", 1, strings.Repeat("a", 250), false},     // Edge case: maximum length
		{"Update with leading/trailing spaces", 1, "   Updated Content   ", false}, // Valid: with leading and trailing white spaces

		// Invalid cases
		{"Post not found", 99, "Updated Content", true},         // Non-existent post ID
		{"Empty content", 1, "", true},                          // Invalid: empty content
		{"Content too long", 1, strings.Repeat("a", 251), true}, // Invalid: content exceeds max length
		{"Update with only white spaces", 1, "     ", true},     // Invalid: only spaces
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			updatedPost, err := UpdatePost(testCase.id, testCase.newContent)

			// Check if error expectation matches
			if (err != nil) != testCase.wantErr {
				t.Fatalf("Expected error: %v, got: %v", testCase.wantErr, err)
			}

			// Verify post content for valid cases
			if !testCase.wantErr && updatedPost.Content != testCase.newContent {
				t.Errorf("Expected content to be '%s', got '%s'", testCase.newContent, updatedPost.Content)
			}

			// Verify no changes were made for invalid cases
			if testCase.wantErr {
				if updatedPost.ID == testCase.id && updatedPost.Content != "Initial Content" {
					t.Errorf("Content should not have been updated for invalid case. Got: '%s'", updatedPost.Content)
				}
			}
		})
	}
}

func TestGetAllPosts(t *testing.T) {

	tests := []struct {
		name    string
		posts   []models.Post
		wantLen int
	}{
		// Case with few posts
		{"Posts exist", []models.Post{
			{ID: 1, Content: "Post 1", Likes: 10, Comments: []models.Comment{}},
			{ID: 2, Content: "Post 2", Likes: 5, Comments: []models.Comment{}},
			{ID: 3, Content: "Post 3", Likes: 1, Comments: []models.Comment{}},
		}, 3},

		// Case with one post
		{"One post exists", []models.Post{
			{ID: 1, Content: "Post 1", Likes: 10, Comments: []models.Comment{}},
		}, 1},

		// Case with no posts
		{"No posts exist", []models.Post{}, 0},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Reset global posts for each test case
			posts = testCase.posts

			// Call GetAllPosts() to get the posts
			result := GetAllPosts()

			// Check if the length of the result matches the expected length
			if len(result) != testCase.wantLen {
				t.Errorf("Expected %d posts, got %d", testCase.wantLen, len(result))
			}
		})
	}
}

func TestLikePost(t *testing.T) {
	// Mock posts
	posts = []models.Post{
		{ID: 1, Content: "Post 1", Likes: 10, Comments: []models.Comment{}},
	}

	tests := []struct {
		name    string
		postID  int
		wantErr bool
	}{
		{"Valid post ID", 1, false},
		{"Invalid post ID", 99, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			likedPost, err := LikePost(testCase.postID)

			if (err != nil) != testCase.wantErr {
				t.Fatalf("Expected error: %v, got: %v", testCase.wantErr, err)
			}

			newLikes := posts[0].Likes + 1
			if !testCase.wantErr && posts[0].Likes != likedPost.Likes {
				t.Errorf("Expected likes to be incremented to 11, got %d", newLikes)
			}
		})
	}
}

func TestGetPostDetailsByID(t *testing.T) {
	posts = []models.Post{
		{ID: 1, Content: "Post 1", Likes: 10, Comments: []models.Comment{{ID: 1, Text: "Nice post"}}},
		{ID: 2, Content: "Post 2", Likes: 5},
	}

	tests := []struct {
		name      string
		postID    int
		wantErr   bool
		wantLikes int
		wantComms int
	}{
		{"Valid post with comments", 1, false, 10, 1},
		{"Valid post without comments", 2, false, 5, 0},
		{"Invalid post ID", 3, true, 0, 0},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			post, err := GetPostDetailsByID(testCase.postID)

			if (err != nil) != testCase.wantErr {
				t.Fatalf("Expected error: %v, got: %v", testCase.wantErr, err)
			}

			if !testCase.wantErr {
				if post.Likes != testCase.wantLikes {
					t.Errorf("Expected likes: %d, got: %d", testCase.wantLikes, post.Likes)
				}

				if len(post.Comments) != testCase.wantComms {
					t.Errorf("Expected %d comments, got %d", testCase.wantComms, len(post.Comments))
				}
			}
		})
	}
}

func TestAddComment(t *testing.T) {
	posts = []models.Post{
		{ID: 1, Content: "Post 1", Likes: 10, Comments: []models.Comment{}},
	}

	tests := []struct {
		name      string
		postID    int
		comment   models.Comment
		wantErr   bool
		wantComms int
	}{
		{"Valid comment for existing post", 1, models.Comment{ID: 1, Text: "Great post!", CreatedAt: time.Now()}, false, 1},
		{"Empty comment for existing post", 1, models.Comment{ID: 2, Text: "", CreatedAt: time.Now()}, true, 0},
		{"Too long comment", 1, models.Comment{ID: 3, Text: strings.Repeat("a", 151), CreatedAt: time.Now()}, true, 0},
		{"Comment for non-existent post", 2, models.Comment{ID: 1, Text: "Interesting!", CreatedAt: time.Now()}, true, 0},
		{"Comment with all white spaces", 2, models.Comment{ID: 1, Text: "     ", CreatedAt: time.Now()}, true, 0},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := AddComment(testCase.postID, testCase.comment)

			// Check for error expectation
			if (err != nil) != testCase.wantErr {
				t.Fatalf("Expected error: %v, got: %v", testCase.wantErr, err)
			}

			if !testCase.wantErr {
				// Check if the comment was added
				post, _ := GetPostDetailsByID(testCase.postID)

				// Check if the comment was added correctly
				if len(post.Comments) != testCase.wantComms {
					t.Errorf("Expected %d comments, got %d", testCase.wantComms, len(post.Comments))
				}

				// Check if the comment text matches
				if post.Comments[len(post.Comments)-1].Text != testCase.comment.Text {
					t.Errorf("Expected comment text '%s', got '%s'", testCase.comment.Text, post.Comments[len(post.Comments)-1].Text)
				}
			}
		})
	}
}
