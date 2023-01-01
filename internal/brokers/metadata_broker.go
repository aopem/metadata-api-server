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

func CreateMetadataBroker() *MetadataBroker {
	storageDirectory := filepath.Join(utils.MainDirectory(), "localStore")
	utils.CreateFolder(storageDirectory)
	return &MetadataBroker{
		storageDirectory: storageDirectory,
	}
}

func (mb *MetadataBroker) CreateMetadata(metadataStore *models.MetadataStore) (*models.MetadataStore, error) {
	// get YAML to write to file
	writeData, err := yaml.Marshal(&metadataStore)
	if err != nil {
		return nil, err
	}

	// get filepath for saving data
	metadataFilepath := filepath.Join(mb.storageDirectory, metadataStore.Id+".yaml")

	// write metadata to file
	if err := utils.WriteFile(metadataFilepath, writeData); err != nil {
		return nil, err
	}

	return metadataStore, nil
}

func (mb *MetadataBroker) DeleteMetadataById(id string) (*models.MetadataStore, error) {
	metadataFilepath := filepath.Join(mb.storageDirectory, id+".yaml")

	// first, get object to return
	metadataStore, err := mb.GetMetadataById(id)
	if err != nil {
		return nil, err
	}

	// then, delete file containing data
	if err := os.Remove(metadataFilepath); err != nil {
		return nil, err
	}

	return metadataStore, nil
}

func (mb *MetadataBroker) GetMetadataById(id string) (*models.MetadataStore, error) {
	metadataFilepath := filepath.Join(mb.storageDirectory, id+".yaml")
	data, err := utils.ReadFile(metadataFilepath)
	if err != nil {
		return nil, err
	}

	// read data into metadata object
	metadataStore := &models.MetadataStore{}
	if err := yaml.Unmarshal(data, metadataStore); err != nil {
		return nil, err
	}

	return metadataStore, nil
}

func (mb *MetadataBroker) GetMetadataList() ([]models.MetadataStore, error) {
	files, err := utils.GetFolderItems(mb.storageDirectory)
	if err != nil {
		return nil, err
	}

	metadataList := []models.MetadataStore{}
	for _, file := range files {
		metadataFilepath := filepath.Join(mb.storageDirectory, file.Name())

		// read data for each file
		data, err := utils.ReadFile(metadataFilepath)
		if err != nil {
			return nil, err
		}

		// then, move to struct
		metadataStore := &models.MetadataStore{}
		if err := yaml.Unmarshal(data, metadataStore); err != nil {
			return nil, err
		}

		// lastly, add to list
		metadataList = append(metadataList, *metadataStore)
	}

	return metadataList, nil
}
