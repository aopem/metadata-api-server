package core

import (
	"metadata-api-server/models"

	"github.com/gin-gonic/gin"
)

type MetadataBroker interface {
	CreateMetadata(bodyData []byte) *models.MetadataStore
	DeleteMetadataById(id string) *models.MetadataStore
	GetMetadataById(id string) *models.MetadataStore
	GetMetadataList() []models.MetadataStore
}

type MetadataService interface {
	CreateMetadata(bodyData []byte) *models.MetadataStore
	DeleteMetadataById(id string) *models.MetadataStore
	GetMetadataById(id string) *models.MetadataStore
	GetMetadata() []models.MetadataStore
}

type MetadataController interface {
	PutMetadata(c *gin.Context)
	DeleteMetadataById(c *gin.Context)
	GetMetadata(c *gin.Context)
	GetMetadataById(c *gin.Context)
}

type Server interface {
	Run(addr string)
}
