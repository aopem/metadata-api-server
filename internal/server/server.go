package server

import (
	"metadata-api-server/internal/brokers"
	"metadata-api-server/internal/controllers"
	"metadata-api-server/internal/query"
	"metadata-api-server/internal/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router             *gin.Engine
	MetadataController *controllers.MetadataController
	QueryController    *controllers.QueryController
}

func CreateServer(router *gin.Engine, mainDirectory string) *Server {
	ib := brokers.CreateIndexBroker(mainDirectory)
	se := query.CreateSearchEngine(ib)

	// create metadata dependencies
	mb := brokers.CreateMetadataBroker(mainDirectory)
	ms := services.CreateMetadataService(mb, se)
	mc := controllers.CreateMetadataController(ms)

	// create query dependencies
	qs := services.CreateQueryService(mb, se)
	qc := controllers.CreateQueryController(qs, ms)

	// create actual server, then route endpoints
	s := &Server{
		Router:             router,
		MetadataController: mc,
		QueryController:    qc,
	}
	s.route()
	return s
}

func (s *Server) Run(addr string) {
	s.Router.Run(addr)
}

func (s *Server) route() {
	// metadata CRUD endpoints
	s.Router.PUT("/metadata", s.MetadataController.PutMetadata)
	s.Router.DELETE("/metadata/:id", s.MetadataController.DeleteMetadataById)
	s.Router.GET("/metadata", s.MetadataController.GetMetadata)
	s.Router.GET("/metadata/:id", s.MetadataController.GetMetadataById)

	// query endpoints
	s.Router.PUT("/metadata/query", s.QueryController.PutMetadataQuery)
}
