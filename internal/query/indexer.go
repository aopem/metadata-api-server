package query

import (
	"fmt"
	"metadata-api-server/models"
	"strings"
)

type Indexer struct {
	cache map[string]map[string][]string
}

func CreateIndexer() *Indexer {
	// initialize internal map
	cache := map[string]map[string][]string{}
	cache["Title"] = map[string][]string{}
	cache["Version"] = map[string][]string{}
	cache["Company"] = map[string][]string{}
	cache["Website"] = map[string][]string{}
	cache["Source"] = map[string][]string{}
	cache["License"] = map[string][]string{}
	cache["Description"] = map[string][]string{}
	cache["Email"] = map[string][]string{}
	cache["Name"] = map[string][]string{}
	return &Indexer{
		cache: cache,
	}
}

func (in *Indexer) IndexMetadata(metadataStore *models.MetadataStore) {
	id := metadataStore.Id
	metadata := metadataStore.Metadata

	// save all metadata in cache
	in.cache["Title"][strings.ToLower(metadata.Title)] = append(in.cache["Title"][strings.ToLower(metadata.Title)], id)
	in.cache["Version"][strings.ToLower(metadata.Version)] = append(in.cache["Version"][strings.ToLower(metadata.Version)], id)
	in.cache["Company"][strings.ToLower(metadata.Company)] = append(in.cache["Company"][strings.ToLower(metadata.Company)], id)
	in.cache["Website"][strings.ToLower(metadata.Website)] = append(in.cache["Website"][strings.ToLower(metadata.Website)], id)
	in.cache["Source"][strings.ToLower(metadata.Source)] = append(in.cache["Source"][strings.ToLower(metadata.Source)], id)
	in.cache["License"][strings.ToLower(metadata.License)] = append(in.cache["License"][strings.ToLower(metadata.License)], id)
	in.cache["Description"][strings.ToLower(metadata.Description)] = append(in.cache["Description"][strings.ToLower(metadata.Description)], id)

	// index all maintainer data
	for _, maintainer := range metadata.Maintainers {
		in.cache["Email"][strings.ToLower(maintainer.Email)] = append(in.cache["Email"][strings.ToLower(maintainer.Email)], id)
		in.cache["Name"][strings.ToLower(maintainer.Name)] = append(in.cache["Name"][strings.ToLower(maintainer.Name)], id)
	}
}

func (in *Indexer) PrintCache() {
	for _, value := range in.cache {
		fmt.Println(value)
	}
}
