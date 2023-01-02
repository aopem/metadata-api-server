package services

import (
	"errors"
	"log"
	"metadata-api-server/internal/core"
	"metadata-api-server/models"
	"strings"
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
	queryMap := qs.QueryToMap(query)

	// search all fields that are non-empty
	// and merge all results into "matches" set
	// which will contain ID of every document that
	// has a partial or full text match to the query
	initialQuery := true
	for field, searchText := range queryMap {
		if searchText == "" {
			continue
		}

		// first search uses "OR" semantics since an "AND" search is limited
		// by the content that is already present in the "matches" passed in
		if initialQuery {
			qs.searchEngine.MetadataFieldOrSearch(field, searchText, matches)
			initialQuery = false
		} else {
			qs.searchEngine.MetadataFieldAndSearch(field, searchText, matches)
		}
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

func (qs *QueryService) QueryToMap(query *models.Query) map[string]string {
	return map[string]string{
		"Title":       strings.ToLower(query.Title),
		"Version":     strings.ToLower(query.Version),
		"Name":        strings.ToLower(query.MaintainerName),
		"Email":       strings.ToLower(query.MaintainerEmail),
		"Company":     strings.ToLower(query.Company),
		"Website":     strings.ToLower(query.Website),
		"Source":      strings.ToLower(query.Source),
		"License":     strings.ToLower(query.License),
		"Description": strings.ToLower(query.Description),
	}
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
