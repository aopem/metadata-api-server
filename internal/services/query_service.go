package services

import (
	"errors"
	"log"
	"metadata-api-server/internal/core"
	"metadata-api-server/models"
)

type QueryService struct {
	MetadataBroker core.MetadataBroker
	searchEngine   core.SearchEngine
}

func CreateQueryService(mb core.MetadataBroker, se core.SearchEngine) *QueryService {
	return &QueryService{
		MetadataBroker: mb,
		searchEngine:   se,
	}
}

func (qs *QueryService) ExecuteQuery(query *models.Query) ([]string, error) {
	if err := qs.validateQuery(query); err != nil {
		return nil, err
	}

	log.Print("Executing query:")
	log.Printf("%+v", *query)

	// uses "matches" as a "set" to store results of a search
	matches := map[string]bool{}

	// search all fields that are non-empty
	// and merge all results into "matches" set
	// which will contain ID of every document that
	// has a partial or full text match to the query
	if query.Title != "" {
		// first search uses "OR" semantics since an "AND" search is limited
		// by the content that is already present in the "matches" passed in
		qs.searchEngine.MetadataFieldOrSearch("Title", query.Title, matches)
	}
	if query.Version != "" {
		qs.searchEngine.MetadataFieldAndSearch("Version", query.Version, matches)
	}
	if query.MaintainerName != "" {
		qs.searchEngine.MetadataFieldAndSearch("Name", query.MaintainerName, matches)
	}
	if query.MaintainerEmail != "" {
		qs.searchEngine.MetadataFieldAndSearch("Email", query.MaintainerEmail, matches)
	}
	if query.Company != "" {
		qs.searchEngine.MetadataFieldAndSearch("Company", query.Company, matches)
	}
	if query.Website != "" {
		qs.searchEngine.MetadataFieldAndSearch("Website", query.Website, matches)
	}
	if query.Source != "" {
		qs.searchEngine.MetadataFieldAndSearch("Source", query.Source, matches)
	}
	if query.License != "" {
		qs.searchEngine.MetadataFieldAndSearch("License", query.License, matches)
	}
	if query.Description != "" {
		qs.searchEngine.MetadataFieldAndSearch("Description", query.Description, matches)
	}

	// convert "matches" to a simple string slice
	matchIds := make([]string, 0)
	for id := range matches {
		matchIds = append(matchIds, id)
	}

	log.Print("Matching IDs found for query:")
	log.Print(matchIds)
	return matchIds, nil
}

func (qs *QueryService) validateQuery(query *models.Query) error {
	if query.Title == "" &&
		query.Version == "" &&
		query.MaintainerName == "" &&
		query.MaintainerEmail == "" &&
		query.Company == "" &&
		query.Website == "" &&
		query.Source == "" &&
		query.License == "" &&
		query.Description == "" {
		return errors.New("All query fields are empty")
	}

	return nil
}
