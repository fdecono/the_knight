package main

import (
	"log"
	"the_knight/internal/web"
)

func main() {
	server := web.NewServer()

	// Start server on port 8080
	// In production, use environment variables or config for port
	if err := server.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
