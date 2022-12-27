package services

import (
	"metadata-api-server/internal/brokers"
	"metadata-api-server/models"
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
	return ms.MetadataBroker.CreateMetadata(bodyData)
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
