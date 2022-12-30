package services

import (
	"metadata-api-server/internal/brokers"
	"metadata-api-server/internal/query"
	"metadata-api-server/models"
)

type QueryService struct {
	MetadataBroker *brokers.MetadataBroker
	searchEngine   *query.SearchEngine
}

func CreateQueryService(mb *brokers.MetadataBroker, se *query.SearchEngine) *QueryService {
	return &QueryService{
		MetadataBroker: mb,
		searchEngine:   se,
	}
}

func (qs *QueryService) ExecuteQuery(query *models.Query) []string {
	// uses "matches" as a "set" to store results of a search
	matches := map[string]bool{}

	// search all fields that are non-empty
	// and merge all results into "matches" set
	// which will contain ID of every document that
	// has a partial or full text match to the query
	if query.Title != "" {
		qs.searchEngine.SearchMetadataField("Title", query.Title, matches)
	}
	if query.Version != "" {
		qs.searchEngine.SearchMetadataField("Version", query.Version, matches)
	}
	if query.MaintainerName != "" {
		qs.searchEngine.SearchMetadataField("Name", query.MaintainerName, matches)
	}
	if query.MaintainerEmail != "" {
		qs.searchEngine.SearchMetadataField("Email", query.MaintainerEmail, matches)
	}
	if query.Company != "" {
		qs.searchEngine.SearchMetadataField("Company", query.Company, matches)
	}
	if query.Website != "" {
		qs.searchEngine.SearchMetadataField("Website", query.Website, matches)
	}
	if query.Source != "" {
		qs.searchEngine.SearchMetadataField("Source", query.Source, matches)
	}
	if query.License != "" {
		qs.searchEngine.SearchMetadataField("License", query.License, matches)
	}
	if query.Description != "" {
		qs.searchEngine.SearchMetadataField("Description", query.Description, matches)
	}

	// convert "matches" to a simple string slice
	matchIds := make([]string, 0)
	for id := range matches {
		matchIds = append(matchIds, id)
	}

	return matchIds
}
