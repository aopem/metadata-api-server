package query

import (
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

func (se *SearchEngine) IndexMetadata(metadataStore *models.MetadataStore) {
	id := metadataStore.Id
	metadata := metadataStore.Metadata
	index := se.IndexBroker.GetIndex()

	// save all metadata in index
	index["Title"][strings.ToLower(metadata.Title)] = append(index["Title"][strings.ToLower(metadata.Title)], id)
	index["Version"][strings.ToLower(metadata.Version)] = append(index["Version"][strings.ToLower(metadata.Version)], id)
	index["Company"][strings.ToLower(metadata.Company)] = append(index["Company"][strings.ToLower(metadata.Company)], id)
	index["Website"][strings.ToLower(metadata.Website)] = append(index["Website"][strings.ToLower(metadata.Website)], id)
	index["Source"][strings.ToLower(metadata.Source)] = append(index["Source"][strings.ToLower(metadata.Source)], id)
	index["License"][strings.ToLower(metadata.License)] = append(index["License"][strings.ToLower(metadata.License)], id)
	index["Description"][strings.ToLower(metadata.Description)] = append(index["Description"][strings.ToLower(metadata.Description)], id)

	// index all maintainer data
	for _, maintainer := range metadata.Maintainers {
		index["Email"][strings.ToLower(maintainer.Email)] = append(index["Email"][strings.ToLower(maintainer.Email)], id)
		index["Name"][strings.ToLower(maintainer.Name)] = append(index["Name"][strings.ToLower(maintainer.Name)], id)
	}
}

func (se *SearchEngine) SearchMetadataField(field string, searchText string, matches map[string]bool) {
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
