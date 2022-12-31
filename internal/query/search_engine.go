package query

import (
	"log"
	"metadata-api-server/internal/brokers"
	"metadata-api-server/models"
	"strings"
)

type SearchEngine struct {
	IndexBroker *brokers.IndexBroker
}

func CreateSearchEngine(ib *brokers.IndexBroker) *SearchEngine {
	return &SearchEngine{
		IndexBroker: ib,
	}
}

func (se *SearchEngine) MetadataFieldOrSearch(field string, searchText string, matches map[string]bool) {
	log.Printf("Searching Metadata field \"%s\" for text \"%s\"", field, searchText)

	// search index[field] for any text that matches
	// the given searchText. Any matches are added to the
	// "matches" set
	for text, ids := range se.IndexBroker.GetIndex()[field] {
		if strings.Contains(text, searchText) {
			for i := range ids {
				matches[ids[i]] = true
			}
		}
	}
}

func (se *SearchEngine) MetadataFieldAndSearch(field string, searchText string, matches map[string]bool) {
	newMatches := map[string]bool{}
	se.MetadataFieldOrSearch(field, searchText, newMatches)

	// now merge "newMatches" with "matches" to make sure that
	// all searches executed with "matches" use "AND" semantics
	// this means any ID from matches that is not also found in
	// "newMatches," then it must be removed from the query results
	deleteIds := make([]string, 0)
	for id := range matches {
		if !newMatches[id] {
			deleteIds = append(deleteIds, id)
		}
	}

	for i := range deleteIds {
		delete(matches, deleteIds[i])
	}
}

func (se *SearchEngine) CreateMetadataIndex(metadataStore *models.MetadataStore) {
	id := metadataStore.Id
	metadata := metadataStore.Metadata

	// save all metadata in index
	log.Print("Indexing Metadata...")
	se.indexField("Title", strings.ToLower(metadata.Title), id)
	se.indexField("Version", strings.ToLower(metadata.Version), id)
	se.indexField("Company", strings.ToLower(metadata.Company), id)
	se.indexField("Website", strings.ToLower(metadata.Website), id)
	se.indexField("Source", strings.ToLower(metadata.Source), id)
	se.indexField("License", strings.ToLower(metadata.License), id)
	se.indexField("Description", strings.ToLower(metadata.Description), id)

	// index all maintainer data
	for _, maintainer := range metadata.Maintainers {
		se.indexField("Email", strings.ToLower(maintainer.Email), id)
		se.indexField("Name", strings.ToLower(maintainer.Name), id)
	}
}

func (se *SearchEngine) DeleteMetadataIndexById(id string) {
	log.Printf("Deleting indexes for Metadata ID \"%s\"...", id)
	index := se.IndexBroker.GetIndex()
	for field, fieldDataMap := range index {
		for fieldData, idList := range fieldDataMap {
			for i := len(idList) - 1; i >= 0; i-- {
				if idList[i] == id {
					index[field][fieldData] = append(
						index[field][fieldData][:i],
						index[field][fieldData][i+1:]...)
				}
			}
		}
	}
}

func (se *SearchEngine) indexField(field string, fieldData string, id string) {
	log.Printf("Indexing field \"%s\" for ID \"%s\"...", field, id)
	index := se.IndexBroker.GetIndex()
	for i := range index[field][fieldData] {
		if index[field][fieldData][i] == id {
			return
		}
	}

	index[field][fieldData] = append(index[field][fieldData], id)
	log.Printf("Added index for Metadata ID \"%s\" successfully", id)
}
