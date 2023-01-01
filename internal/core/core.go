package core

import (
	"metadata-api-server/models"

	"github.com/gin-gonic/gin"
)

type MetadataBroker interface {
	CreateMetadata(metadataStore *models.MetadataStore) (*models.MetadataStore, error)
	DeleteMetadataById(id string) (*models.MetadataStore, error)
	GetMetadataById(id string) (*models.MetadataStore, error)
	GetMetadataList() ([]models.MetadataStore, error)
	GetStorageDirectory() string
}

type MetadataService interface {
	CreateMetadata(metadata *models.Metadata) (*models.MetadataStore, error)
	DeleteMetadataById(id string) (*models.MetadataStore, error)
	GetMetadataById(id string) (*models.MetadataStore, error)
	GetMetadata() ([]models.MetadataStore, error)
}

type MetadataController interface {
	PutMetadata(c *gin.Context)
	DeleteMetadataById(c *gin.Context)
	GetMetadata(c *gin.Context)
	GetMetadataById(c *gin.Context)
}

type QueryService interface {
	ExecuteQuery(query *models.Query) ([]string, error)
}

type QueryController interface {
	PutMetadataQuery(c *gin.Context)
}

type IndexBroker interface {
	CreateIndex(metadataStore *models.MetadataStore)
	DeleteIndexById(id string)
	GetIndex() map[string]map[string][]string
	GetIndexPath() string
	SaveIndex() error
	IndexEmpty() bool
}

type SearchEngine interface {
	MetadataFieldOrSearch(field string, searchText string, matches map[string]bool)
	MetadataFieldAndSearch(field string, searchText string, matches map[string]bool)
}

type Server interface {
	Run(addr string)
}
