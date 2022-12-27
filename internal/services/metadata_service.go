package services

import (
	"metadata-api-server/internal/brokers"
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"

	"gopkg.in/yaml.v2"
)

type MetadataService struct {
	MetadataBroker *brokers.MetadataBroker
}

func CreateMetadataService(b *brokers.MetadataBroker) *MetadataService {
	return &MetadataService{
		MetadataBroker: b,
	}
}

func (ms *MetadataService) CreateMetadata(bodyData []byte) *models.MetadataStore {
	metadata := &models.Metadata{}
	err := yaml.Unmarshal(bodyData, metadata)
	if err != nil {
		return nil
	}

	// get hash value and add to given metadata info
	metadataStore := &models.MetadataStore{
		Id:       utils.CalculateHash(bodyData),
		Metadata: metadata,
	}

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
