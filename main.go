package main

import (
	"log"
	"metadata-api-server/internal/server"
	"metadata-api-server/internal/utils"
	"path/filepath"
)

func main() {
	addr := "localhost:8080"

	log.Printf("Initializing server at \"%s\"...", addr)
	rootDirectory := utils.MainDirectory()
	s := server.CreateServer(
		filepath.Join(rootDirectory, "localIndex"),
		filepath.Join(rootDirectory, "localStore"))
	s.Run(addr)
}
