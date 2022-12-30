package brokers

import (
	"encoding/gob"
	"metadata-api-server/internal/utils"
	"os"
	"path/filepath"
)

type IndexBroker struct {
	indexData map[string]map[string][]string
	indexPath string
	indexFile *os.File
}

func CreateIndexBroker(mainDirectory string) *IndexBroker {
	indexPath := filepath.Join(mainDirectory, "localIndex", "index.gob")
	indexFile := &os.File{}
	indexData := map[string]map[string][]string{}

	// if index file does not already exist,
	// then initialize all necessary attributes
	if !utils.FileExists(indexPath) || utils.FileEmpty(indexPath) {
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
		if err := os.MkdirAll(filepath.Dir(indexPath), os.ModePerm); err != nil {
			// throw error
			return nil
		}
		indexFile = utils.CreateFile(indexPath)
	} else {
		// otherwise, decode existing file data
		decoder := gob.NewDecoder(indexFile)
		decoder.Decode(&indexData)
	}

	return &IndexBroker{
		indexData: indexData,
		indexPath: indexPath,
		indexFile: indexFile,
	}
}

func (ib *IndexBroker) GetIndex() map[string]map[string][]string {
	return ib.indexData
}

func (ib *IndexBroker) SaveIndex() {
	encoder := gob.NewEncoder(ib.indexFile)
	encoder.Encode(ib.indexData)
}