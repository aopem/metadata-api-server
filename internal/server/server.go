package server

import (
	"metadata-api-server/internal/brokers"
	"metadata-api-server/internal/controllers"
	"metadata-api-server/internal/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router             *gin.Engine
	MetadataController *controllers.MetadataController
}

func CreateServer(router *gin.Engine, mainDirectory string) *Server {
	mb := brokers.CreateMetadataBroker(mainDirectory)
	ms := services.CreateMetadataService(mb)
	mc := controllers.CreateMetadataController(ms)

	s := &Server{
		Router:             router,
		MetadataController: mc,
	}

	s.route()
	return s
}

func (s *Server) Run(addr string) {
	s.Router.Run(addr)
}

func (s *Server) route() {
	s.Router.PUT("/metadata", s.MetadataController.PutMetadata)
	s.Router.DELETE("/metadata/:id", s.MetadataController.DeleteMetadataById)
	s.Router.GET("/metadata", s.MetadataController.GetMetadata)
	s.Router.GET("/metadata/:id", s.MetadataController.GetMetadataById)
}
