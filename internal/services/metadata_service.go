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
	indexer        *query.Indexer
}

func CreateMetadataService(b *brokers.MetadataBroker) *MetadataService {
	return &MetadataService{
		MetadataBroker: b,
		indexer:        query.CreateIndexer(),
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
	ms.indexer.IndexMetadata(metadataStore)
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
