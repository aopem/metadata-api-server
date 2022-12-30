package query

import (
	"metadata-api-server/models"
	"strings"
)

type SearchEngine struct {
	index map[string]map[string][]string
}

func CreateSearchEngine() *SearchEngine {
	// initialize internal map
	index := map[string]map[string][]string{}
	index["Title"] = map[string][]string{}
	index["Version"] = map[string][]string{}
	index["Company"] = map[string][]string{}
	index["Website"] = map[string][]string{}
	index["Source"] = map[string][]string{}
	index["License"] = map[string][]string{}
	index["Description"] = map[string][]string{}
	index["Email"] = map[string][]string{}
	index["Name"] = map[string][]string{}
	return &SearchEngine{
		index: index,
	}
}

func (se *SearchEngine) IndexMetadata(metadataStore *models.MetadataStore) {
	id := metadataStore.Id
	metadata := metadataStore.Metadata

	// save all metadata in index
	se.index["Title"][strings.ToLower(metadata.Title)] = append(se.index["Title"][strings.ToLower(metadata.Title)], id)
	se.index["Version"][strings.ToLower(metadata.Version)] = append(se.index["Version"][strings.ToLower(metadata.Version)], id)
	se.index["Company"][strings.ToLower(metadata.Company)] = append(se.index["Company"][strings.ToLower(metadata.Company)], id)
	se.index["Website"][strings.ToLower(metadata.Website)] = append(se.index["Website"][strings.ToLower(metadata.Website)], id)
	se.index["Source"][strings.ToLower(metadata.Source)] = append(se.index["Source"][strings.ToLower(metadata.Source)], id)
	se.index["License"][strings.ToLower(metadata.License)] = append(se.index["License"][strings.ToLower(metadata.License)], id)
	se.index["Description"][strings.ToLower(metadata.Description)] = append(se.index["Description"][strings.ToLower(metadata.Description)], id)

	// index all maintainer data
	for _, maintainer := range metadata.Maintainers {
		se.index["Email"][strings.ToLower(maintainer.Email)] = append(se.index["Email"][strings.ToLower(maintainer.Email)], id)
		se.index["Name"][strings.ToLower(maintainer.Name)] = append(se.index["Name"][strings.ToLower(maintainer.Name)], id)
	}
}

func (se *SearchEngine) SearchMetadataField(field string, searchText string) map[string]bool {
	// create a "set" to store all metadata IDs that match
	matches := make(map[string]bool)

	// search se.index[field] for any text that matches
	// the given searchText. Any matches are added to the
	// "matches" set
	for text, ids := range se.index[field] {
		if strings.Contains(text, searchText) {
			for i := range ids {
				matches[ids[i]] = true
			}
		}
	}

	return matches
}
