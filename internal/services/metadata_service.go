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
	indexBroker    core.IndexBroker
}

func CreateMetadataService(mb core.MetadataBroker, ib core.IndexBroker) *MetadataService {
	return &MetadataService{
		MetadataBroker: mb,
		indexBroker:    ib,
	}
}

func (ms *MetadataService) CreateMetadata(metadataStore *models.MetadataStore) (*models.MetadataStore, error) {
	// check if id is empty - this means a new record is being created
	if metadataStore.Id == "" {
		// get metadata in byte format to calculate a unique hash
		metadataBytes, err := yaml.Marshal(*metadataStore.Metadata)
		if err != nil {
			return nil, err
		}

		// add hash ID to metadata for storage
		hash, err := utils.CalculateHash(metadataBytes)
		if err != nil {
			return nil, err
		}

		metadataStore.Id = hash

		// index new metadata for searches, then create record
		log.Printf("Creating Metadata ID \"%s\"...", metadataStore.Id)
		ms.indexBroker.CreateIndex(metadataStore)
		return ms.MetadataBroker.CreateMetadata(metadataStore)
	}

	log.Printf("Metadata ID \"%s\" already exists, retrieving existing data...", metadataStore.Id)
	metadataExisting, err := ms.GetMetadataById(metadataStore.Id)
	if err != nil {
		return nil, err
	}

	// if ID matches, but this metadata has new content, then update existing record
	if !utils.MetadataEqual(metadataStore, metadataExisting) {
		log.Printf("Changes detected. Updating Metadata ID \"%s\" with new data...", metadataStore.Id)
		return ms.MetadataBroker.CreateMetadata(metadataStore)
	}

	return metadataExisting, nil
}

func (ms *MetadataService) DeleteMetadataById(id string) (*models.MetadataStore, error) {
	if err := ms.validateId(id); err != nil {
		return nil, err
	}

	// delete from index, then delete from local store
	log.Printf("Deleting Metadata ID \"%s\"...", id)
	ms.indexBroker.DeleteIndexById(id)
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
