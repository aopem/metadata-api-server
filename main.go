package main

import (
	"metadata-api-server/internal/server"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
)

func mainDirectory() string {
	_, filename, _, ok := runtime.Caller(1)

	if !ok {
		// TODO: throw error
		return ""
	}

	return filepath.Dir(filename)
}

func main() {
	s := server.CreateServer(gin.Default(), mainDirectory())
	s.Run("localhost:8080")
}
