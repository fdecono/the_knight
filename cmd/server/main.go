package main

import (
	"log"
	"os"
	"the_knight/internal/web"
)

func main() {
	server := web.NewServer()

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	addr := ":" + port
	log.Printf("Starting Knight's Tour server on %s", addr)
	if err := server.Start(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
