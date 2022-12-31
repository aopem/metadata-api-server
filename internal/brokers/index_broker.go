package brokers

import (
	"encoding/gob"
	"log"
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
		if err := os.MkdirAll(filepath.Dir(indexPath), os.ModePerm); err != nil {
			// throw error
			return nil
		}
		utils.CreateFile(indexPath)
	} else {
		log.Print("Local index already exists, loading...")

		indexFile, err := os.Open(indexPath)
		if err != nil {
			panic(err)
		}

		// otherwise, decode existing file data
		decoder := gob.NewDecoder(indexFile)
		decoder.Decode(&indexData)
	}

	indexFile, err := os.OpenFile(indexPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
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
	log.Printf("Saving local index at %s...", ib.indexPath)
	encoder := gob.NewEncoder(ib.indexFile)
	err := encoder.Encode(ib.indexData)
	if err != nil {
		panic(err)
	}
}
