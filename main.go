package main

import (
	"log"
	"metadata-api-server/internal/server"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
)

func mainDirectory() string {
	_, filename, _, ok := runtime.Caller(1)

	if !ok {
		log.Fatal("Error: Could not retrieve main.go directory")
	}

	return filepath.Dir(filename)
}

func main() {
	mainDirectory := mainDirectory()
	addr := "localhost:8080"

	log.Printf("Program main.go is located at %s", mainDirectory)
	log.Printf("Initializing server at \"%s\"...", addr)
	s := server.CreateServer(gin.Default(), mainDirectory)
	s.Run(addr)
}
