package main

// This is the legacy CLI entry point.
// The new web server entry point is in cmd/server/main.go
// Run with: go run cmd/server/main.go

import (
	"fmt"
	"the_knight/internal/web"
)

func main() {
	fmt.Println("Starting Knight's Tour Web Server...")
	fmt.Println("Visit http://localhost:8080 in your browser")

	server := web.NewServer()
	if err := server.Start(":8080"); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
