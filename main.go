package main

import (
	"mini-social-media-api/routes"
)

func main() {
	// Initialize routes and start the HTTP server on port 8081
	router := routes.InitRoutes()
	router.Run(":8081")		
}