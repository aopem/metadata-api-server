package query

import (
	"log"
	"metadata-api-server/internal/core"
	"strings"
)

type SearchEngine struct {
	IndexBroker core.IndexBroker
}

func CreateSearchEngine(ib core.IndexBroker) *SearchEngine {
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
