package main

import (
	"log"
	"net/http"
	"os"

	"devsite/internal/content"
	"devsite/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Load content from markdown files
	loader := content.NewLoader("content")
	if err := loader.Load(); err != nil {
		log.Fatalf("failed to load content: %v", err)
	}

	// Create and start server
	srv := server.New(loader)

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, srv.Router()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
