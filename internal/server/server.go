package server

import (
	"log"
	"metadata-api-server/internal/brokers"
	"metadata-api-server/internal/controllers"
	"metadata-api-server/internal/core"
	"metadata-api-server/internal/query"
	"metadata-api-server/internal/services"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router             *gin.Engine
	MetadataController core.MetadataController
	QueryController    core.QueryController
	metadataService    core.MetadataService
	indexBroker        core.IndexBroker
}

func CreateServer(indexDirectory string, storageDirectory string) *Server {
	ib := brokers.CreateIndexBroker(indexDirectory)
	se := query.CreateSearchEngine(ib)

	// create metadata dependencies
	mb := brokers.CreateMetadataBroker(storageDirectory)
	ms := services.CreateMetadataService(mb, ib)
	mc := controllers.CreateMetadataController(ms)

	// create query dependencies
	qs := services.CreateQueryService(mb, se)
	qc := controllers.CreateQueryController(qs, ms)

	// create actual server, then route endpoints
	s := &Server{
		Router:             gin.Default(),
		MetadataController: mc,
		QueryController:    qc,
		metadataService:    ms,
		indexBroker:        ib,
	}

	log.Print("Routing server endpoints...")
	s.route()
	return s
}

func (s *Server) Run(addr string) {
	server := endless.NewServer(addr, s.Router)

	// load index on startup
	s.onStartup()

	// for saving index on a shutdown signal
	server.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT] = append(
		server.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT],
		s.onShutdown)
	server.SignalHooks[endless.PRE_SIGNAL][syscall.SIGTERM] = append(
		server.SignalHooks[endless.PRE_SIGNAL][syscall.SIGTERM],
		s.onShutdown)

	server.ListenAndServe()
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

func (s *Server) onStartup() {
	log.Print("Checking if stored data needs to be indexed...")

	// see if any metadata is stored locally
	metadataStoreList, err := s.metadataService.GetMetadata()
	if err != nil {
		log.Printf("[FATAL ERROR] %s", err.Error())
	}

	// if index is empty, but stored metadata is not,
	// then load all of stored metadata into index
	if s.indexBroker.IndexEmpty() && len(metadataStoreList) > 0 {
		log.Print("Index empty and locally stored data found. Loading into index...")
		for i := range metadataStoreList {
			s.indexBroker.CreateIndex(&metadataStoreList[i])
		}
	} else {
		log.Print("No pre-indexing required. Skipping...")
	}
}

func (s *Server) onShutdown() {
	log.Print("Saving index data...")
	if err := s.indexBroker.SaveIndex(); err != nil {
		log.Print("[FATAL ERROR] Could not save index on server shutdown")
		log.Printf("[FATAL ERROR] %s", err.Error())
	}
}
