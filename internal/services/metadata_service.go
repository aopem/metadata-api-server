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

func (s *MetadataService) CreateMetadata(bodyData []byte) *models.MetadataStore {
	return s.MetadataBroker.CreateMetadata(bodyData)
}

func (s *MetadataService) GetMetadataByHash(hash string) *models.MetadataStore {
	return s.MetadataBroker.GetMetadataYamlByHash(hash)
}

func (s *MetadataService) GetMetadata() []models.Metadata {
	return s.MetadataBroker.GetMetadataYamlList()
}
