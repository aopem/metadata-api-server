package brokers

import (
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type MetadataBroker struct {
	StorageDirectory string
}

func CreateMetadataBroker(mainDirectory string) *MetadataBroker {
	storageDirectory := filepath.Join(mainDirectory, "localStore")

	if err := os.MkdirAll(storageDirectory, os.ModePerm); err != nil {
		// throw error
		return nil
	}

	return &MetadataBroker{
		StorageDirectory: storageDirectory,
	}
}

func (mb *MetadataBroker) CreateMetadata(bodyData []byte) *models.MetadataStore {
	metadata := &models.Metadata{}
	err := yaml.Unmarshal(bodyData, metadata)
	if err != nil {
		return nil
	}

	// get checksum hash value
	metadataStore := &models.MetadataStore{
		Id:       utils.CalculateHash(bodyData),
		Metadata: metadata,
	}

	// get YAML for writing to file
	writeData, err := yaml.Marshal(&metadataStore)
	if err != nil {
		return nil
	}

	// get filepath for saving data
	metadataFilepath := filepath.Join(mb.StorageDirectory, metadataStore.Id+".yaml")

	// write metadata to file
	utils.WriteFile(metadataFilepath, writeData)
	return metadataStore
}

func (mb *MetadataBroker) GetMetadataYamlById(id string) *models.MetadataStore {
	metadataFilepath := filepath.Join(mb.StorageDirectory, id+".yaml")
	data := utils.ReadFile(metadataFilepath)

	// read data into metadata object
	metadataStore := &models.MetadataStore{}
	err := yaml.Unmarshal(data, metadataStore)
	if err != nil {
		return nil
	}

	return metadataStore
}

func (mb *MetadataBroker) GetMetadataYamlList() []models.Metadata {
	return nil
}
