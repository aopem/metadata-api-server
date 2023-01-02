package brokers

import (
	"encoding/gob"
	"log"
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"
	"os"
	"path/filepath"
	"strings"
)

type IndexBroker struct {
	indexData    map[string]map[string][]string
	indexPath    string
	indexedIdSet map[string]bool
}

func CreateIndexBroker(rootDirectory string) *IndexBroker {
	indexPath := filepath.Join(rootDirectory, "index.gob")
	log.Printf("Local index will be stored at %s", indexPath)

	// initialize index map
	indexData := map[string]map[string][]string{}

	// if index file does not already exist,
	// then initialize all necessary attributes
	if !utils.FileExists(indexPath) || utils.FileEmpty(indexPath) {
		log.Print("Local index either does not exist or is empty")

		// initialize index map
		indexData["Title"] = map[string][]string{}
		indexData["Version"] = map[string][]string{}
		indexData["Company"] = map[string][]string{}
		indexData["Website"] = map[string][]string{}
		indexData["Source"] = map[string][]string{}
		indexData["License"] = map[string][]string{}
		indexData["Description"] = map[string][]string{}
		indexData["Email"] = map[string][]string{}
		indexData["Name"] = map[string][]string{}

		// create a local store for the index
		log.Print("Creating new local index...")
		utils.CreateFolder(filepath.Dir(indexPath))
		utils.CreateFile(indexPath)
	} else {
		// decode existing file data
		log.Print("Local index already exists, loading...")
		indexFile := utils.OpenFile(indexPath, 0, os.ModePerm)
		decoder := gob.NewDecoder(indexFile)
		decoder.Decode(&indexData)
	}

	return &IndexBroker{
		indexData:    indexData,
		indexPath:    indexPath,
		indexedIdSet: make(map[string]bool),
	}
}

func (ib *IndexBroker) CreateIndex(metadataStore *models.MetadataStore) {
	id := metadataStore.Id
	metadata := metadataStore.Metadata

	// save all metadata in index
	log.Print("Indexing Metadata...")
	ib.indexField("Title", strings.ToLower(metadata.Title), id)
	ib.indexField("Version", strings.ToLower(metadata.Version), id)
	ib.indexField("Company", strings.ToLower(metadata.Company), id)
	ib.indexField("Website", strings.ToLower(metadata.Website), id)
	ib.indexField("Source", strings.ToLower(metadata.Source), id)
	ib.indexField("License", strings.ToLower(metadata.License), id)
	ib.indexField("Description", strings.ToLower(metadata.Description), id)

	// index all maintainer data
	for _, maintainer := range metadata.Maintainers {
		ib.indexField("Email", strings.ToLower(maintainer.Email), id)
		ib.indexField("Name", strings.ToLower(maintainer.Name), id)
	}

	// add to indexed ID set
	ib.indexedIdSet[id] = true
}

func (ib *IndexBroker) DeleteIndexById(id string) {
	log.Printf("Deleting indexes for Metadata ID \"%s\"...", id)
	for field, fieldDataMap := range ib.indexData {
		for fieldData, idList := range fieldDataMap {
			for i := len(idList) - 1; i >= 0; i-- {
				if idList[i] == id {
					ib.indexData[field][fieldData] = append(
						ib.indexData[field][fieldData][:i],
						ib.indexData[field][fieldData][i+1:]...)

					// delete from indexed ID set
					delete(ib.indexedIdSet, id)
				}
			}
		}
	}
}

func (ib *IndexBroker) GetIndex() map[string]map[string][]string {
	return ib.indexData
}

func (ib *IndexBroker) GetIndexPath() string {
	return ib.indexPath
}

func (ib *IndexBroker) SaveIndex() error {
	log.Printf("Saving local index at %s...", ib.indexPath)

	// open file, then encode/save ib.indexData to the file
	indexFile := utils.OpenFile(ib.indexPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	encoder := gob.NewEncoder(indexFile)
	if err := encoder.Encode(ib.indexData); err != nil {
		return err
	}

	return nil
}

func (ib *IndexBroker) IndexEmpty() bool {
	return len(ib.indexedIdSet) == 0
}

func (ib *IndexBroker) indexField(field string, fieldData string, id string) {
	log.Printf("Indexing field \"%s\" for ID \"%s\"...", field, id)
	for i := range ib.indexData[field][fieldData] {
		if ib.indexData[field][fieldData][i] == id {
			return
		}
	}

	ib.indexData[field][fieldData] = append(ib.indexData[field][fieldData], id)
	log.Printf("Added index for Metadata ID \"%s\" successfully", id)
}
