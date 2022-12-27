package core

import (
	"metadata-api-server/models"

	"github.com/gin-gonic/gin"
)

type MetadataBroker interface {
	GetMetadataYamlById(id string) *models.Metadata
	GetMetadataYamlList() []models.Metadata
}

type MetadataService interface {
	GetMetadataById(id string) *models.Metadata
	GetMetadata() []models.Metadata
}

type MetadataController interface {
	PutMetadata(c *gin.Context)
	GetMetadata(c *gin.Context)
	GetMetadataById(c *gin.Context)
}

type Server interface {
	Run(addr string)
}
