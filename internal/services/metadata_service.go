package services

import (
	"errors"
	"fmt"
	"log"
	"metadata-api-server/internal/core"
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"
	"regexp"

	"gopkg.in/yaml.v2"
)

type MetadataService struct {
	MetadataBroker core.MetadataBroker
	searchEngine   core.SearchEngine
}

func CreateMetadataService(mb core.MetadataBroker, se core.SearchEngine) *MetadataService {
	return &MetadataService{
		MetadataBroker: mb,
		searchEngine:   se,
	}
}

func (ms *MetadataService) CreateMetadata(metadata *models.Metadata) (*models.MetadataStore, error) {
	// get metadata in byte format to calculate a unique hash
	metadataBytes, err := yaml.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	// add hash ID to metadata for storage
	hash, err := utils.CalculateHash(metadataBytes)
	if err != nil {
		return nil, err
	}

	metadataStore := &models.MetadataStore{
		Id:       hash,
		Metadata: metadata,
	}

	// pre-process for searches, then create using broker
	ms.searchEngine.CreateMetadataIndex(metadataStore)

	log.Printf("Creating Metadata ID \"%s\"...", metadataStore.Id)
	return ms.MetadataBroker.CreateMetadata(metadataStore)
}

func (ms *MetadataService) DeleteMetadataById(id string) (*models.MetadataStore, error) {
	if err := ms.validateId(id); err != nil {
		return nil, err
	}

	// delete from index, then delete from local store
	log.Printf("Deleting Metadata ID \"%s\"...", id)
	ms.searchEngine.DeleteMetadataIndexById(id)
	return ms.MetadataBroker.DeleteMetadataById(id)
}

func (ms *MetadataService) GetMetadataById(id string) (*models.MetadataStore, error) {
	if err := ms.validateId(id); err != nil {
		return nil, err
	}

	log.Printf("Retrieving Metadata ID \"%s\"...", id)
	return ms.MetadataBroker.GetMetadataById(id)
}

func (ms *MetadataService) GetMetadata() ([]models.MetadataStore, error) {
	log.Print("Retrieving all Metadata...")
	return ms.MetadataBroker.GetMetadataList()
}

func (ms *MetadataService) validateId(id string) error {
	regexp.MustCompile("[a-zA-Z0-9]")
	if len(id) != 8 {
		return errors.New(fmt.Sprintf("Invalid ID \"%s\" - must be exactly 8 alphanumeric characters", id))
	}

	return nil
}
