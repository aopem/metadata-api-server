package brokers

import (
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type MetadataBroker struct {
	storageDirectory string
}

func CreateMetadataBroker(mainDirectory string) *MetadataBroker {
	storageDirectory := filepath.Join(mainDirectory, "localStore")

	if err := os.MkdirAll(storageDirectory, os.ModePerm); err != nil {
		// throw error
		return nil
	}

	return &MetadataBroker{
		storageDirectory: storageDirectory,
	}
}

func (mb *MetadataBroker) CreateMetadata(metadataStore *models.MetadataStore) *models.MetadataStore {
	// get YAML to write to file
	writeData, err := yaml.Marshal(&metadataStore)
	if err != nil {
		return nil
	}

	// get filepath for saving data
	metadataFilepath := filepath.Join(mb.storageDirectory, metadataStore.Id+".yaml")

	// write metadata to file
	utils.WriteFile(metadataFilepath, writeData)
	return metadataStore
}

func (mb *MetadataBroker) DeleteMetadataById(id string) *models.MetadataStore {
	metadataFilepath := filepath.Join(mb.storageDirectory, id+".yaml")

	// first, get object to return
	metadataStore := mb.GetMetadataById(id)

	// then, delete file containing data
	err := os.Remove(metadataFilepath)
	if err != nil {
		return nil
	}

	return metadataStore
}

func (mb *MetadataBroker) GetMetadataById(id string) *models.MetadataStore {
	metadataFilepath := filepath.Join(mb.storageDirectory, id+".yaml")
	data := utils.ReadFile(metadataFilepath)

	// read data into metadata object
	metadataStore := &models.MetadataStore{}
	err := yaml.Unmarshal(data, metadataStore)
	if err != nil {
		return nil
	}

	return metadataStore
}

func (mb *MetadataBroker) GetMetadataList() []models.MetadataStore {
	files := utils.GetFolderItems(mb.storageDirectory)

	metadataList := []models.MetadataStore{}
	for _, file := range files {
		metadataFilepath := filepath.Join(mb.storageDirectory, file.Name())

		// read data for each file
		data := utils.ReadFile(metadataFilepath)

		// then, move to struct
		metadataStore := &models.MetadataStore{}
		err := yaml.Unmarshal(data, metadataStore)
		if err != nil {
			return nil
		}

		// lastly, add to list
		metadataList = append(metadataList, *metadataStore)
	}

	return metadataList
}
