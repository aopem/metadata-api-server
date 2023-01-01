package main

import (
	"log"
	"metadata-api-server/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	addr := "localhost:8080"

	log.Printf("Initializing server at \"%s\"...", addr)
	s := server.CreateServer(gin.Default())
	s.Run(addr)
}
