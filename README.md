# mini-social-media-api
The Mini Social Media API is a RESTful API that enables users to create, update, and interact with posts. This solution is designed to simulate core social media functionalities like adding posts, liking posts, and adding comments to posts. The API supports concurrency using mutexes to ensure thread-safe operations and maintain data integrity in a concurrent environment. Additionally, it provides simple pagination for retrieving posts, allowing users to fetch posts in manageable chunks using page and limit query parameters.

Packages used: 
- Gin: Selected for its performance and ease of use, with efficient routing and JSON handling.
- Logrus: Used for structured logging to simplify debugging and track API's activity.

---

## Features
- Create new posts (text-based)
- Update existing posts
- Retrieve all posts
- Like specific posts
- Add comments to specific posts
- Retrieve post details, including comments and likes

---

## Build steps in the local machine
- Clone the repository.
- Install dependencies: ```go mod tidy```
- Run the application: ```go run main.go```
- Access the API: http://localhost:8081

---

## Assumptions
- Data is temporarily stored in memory using Go structs and slices, meaning all data will be lost upon application restart.
- The API does not include user authentication or authorization, assuming all requests are made by authenticated users.
- Concurrency is managed with locking mechanisms (e.g., sync.Mutex) to ensure thread-safe operations on posts.
- Posts are simple text messages without additional attributes like images or user information.
- "Like" functionality increases the like count without distinguishing between unique or repeated likes.
- Updating a post only modifies its content; associated comments and likes remain unaffected.
- Deletion of posts, comments, or likes is not required and is not implemented.

---

## Suggested Improvements
- Use a database (e.g., MongoDB) for data persistence to avoid data loss on restart.
- Add support for media attachments in posts.
- Introduce API rate limiting to prevent misuse.
- Expand unit tests to cover edge cases and add integration tests for end-to-end validation.
- Enable editing and deletion for posts, comments, and likes.
- Add comment threads by allowing replies to specific comments.