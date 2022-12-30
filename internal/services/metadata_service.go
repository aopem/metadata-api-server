package services

import (
	"metadata-api-server/internal/brokers"
	"metadata-api-server/internal/query"
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"

	"gopkg.in/yaml.v2"
)

type MetadataService struct {
	MetadataBroker *brokers.MetadataBroker
	searchEngine   *query.SearchEngine
}

func CreateMetadataService(mb *brokers.MetadataBroker, se *query.SearchEngine) *MetadataService {
	return &MetadataService{
		MetadataBroker: mb,
		searchEngine:   se,
	}
}

func (ms *MetadataService) CreateMetadata(metadata *models.Metadata) *models.MetadataStore {
	// get metadata in byte format to calculate a unique hash
	metadataBytes, err := yaml.Marshal(metadata)
	if err != nil {
		panic(err)
	}

	// add hash ID to metadata for storage
	metadataStore := &models.MetadataStore{
		Id:       utils.CalculateHash(metadataBytes),
		Metadata: metadata,
	}

	// pre-process for searches, then create using broker
	ms.searchEngine.IndexMetadata(metadataStore)
	return ms.MetadataBroker.CreateMetadata(metadataStore)
}

func (ms *MetadataService) DeleteMetadataById(id string) *models.MetadataStore {
	return ms.MetadataBroker.DeleteMetadataById(id)
}

func (ms *MetadataService) GetMetadataById(id string) *models.MetadataStore {
	return ms.MetadataBroker.GetMetadataById(id)
}

func (ms *MetadataService) GetMetadata() []models.MetadataStore {
	return ms.MetadataBroker.GetMetadataList()
}
