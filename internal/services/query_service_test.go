package services

import (
	"metadata-api-server/internal/brokers"
	"metadata-api-server/internal/query"
	"metadata-api-server/internal/testutils"
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryService(t *testing.T) {
	testcases := []testutils.Test{{
		Name:     "TestExecuteQuery",
		Function: TestExecuteQuery,
	}}

	// seed all random numbers that are generated
	testutils.SeedRandomGenerator()

	// run tests
	for i := range testcases {
		// if folders already exist, clean before running tests
		utils.DeleteFolder(testutils.TestStorageDirectory)
		utils.DeleteFolder(testutils.TestIndexDirectory)
		t.Run(testcases[i].Name, testcases[i].Function)
	}

	// cleanup
	utils.DeleteFolder(testutils.TestStorageDirectory)
	utils.DeleteFolder(testutils.TestIndexDirectory)
}

func TestExecuteQuery(t *testing.T) {
	assert := assert.New(t)

	// create dependencies
	ib := brokers.CreateIndexBroker(testutils.TestIndexDirectory)
	se := query.CreateSearchEngine(ib)
	mb := brokers.CreateMetadataBroker(testutils.TestStorageDirectory)
	ms := CreateMetadataService(mb, ib)
	qs := CreateQueryService(mb, se)

	// add some entries before executing a query
	entry1 := &models.MetadataStore{
		Id: "",
		Metadata: &models.Metadata{
			Title:   "Valid App 1",
			Version: "0.0.1",
			Maintainers: []models.Maintainer{{
				Name:  "firstmaintainer app1",
				Email: "firstmaintainer@hotmail.com",
			}, {
				Name:  "secondmaintainer app1",
				Email: "secondmaintainer@gmail.com",
			}},
			Company: "Random Inc.",
			Website: "https://website.com",
			Source:  "https://github.com/random/repo",
			License: "Apache-2.0",
			Description: `|
				### Interesting Title
				Some application content, and description`,
		},
	}

	entry2 := &models.MetadataStore{
		Id: "",
		Metadata: &models.Metadata{
			Title:   "Valid App 2",
			Version: "1.0.1",
			Maintainers: []models.Maintainer{{
				Name:  "AppTwo Maintainer",
				Email: "apptwo@hotmail.com",
			}},
			Company: "Upbound Inc.",
			Website: "https://upbound.io",
			Source:  "https://github.com/upbound/repo",
			License: "Apache-2.0",
			Description: `|
				### Why app 2 is the best
				Because it simply is...`,
		},
	}

	entry1Created, err := ms.CreateMetadata(entry1)
	assert.NoError(err)
	entry2Created, err := ms.CreateMetadata(entry2)
	assert.NoError(err)

	// execute a search for entry1
	matchIds, err := qs.ExecuteQuery(&models.Query{
		Title:           "Valid App 1",
		Version:         "",
		MaintainerName:  "",
		MaintainerEmail: "",
		Company:         "Random Inc.",
		Website:         "https://website.com",
		Source:          "https://github.com/random/repo",
		License:         "",
		Description:     "",
	})
	assert.NoError(err)
	assert.Len(matchIds, 1)
	assert.Contains(matchIds, entry1Created.Id)

	// execute a search for entry2
	matchIds, err = qs.ExecuteQuery(&models.Query{
		Title:           "app",
		Version:         "",
		MaintainerName:  "AppTwo Maintainer",
		MaintainerEmail: "@hotmail.com",
		Company:         "",
		Website:         "",
		Source:          "",
		License:         "",
		Description:     "### Why app 2 is the best",
	})
	assert.NoError(err)
	assert.Len(matchIds, 1)
	assert.Contains(matchIds, entry2Created.Id)

	// execute a search for both entries
	matchIds, err = qs.ExecuteQuery(&models.Query{
		Title:           "app",
		Version:         "",
		MaintainerName:  "maintainer",
		MaintainerEmail: "",
		Company:         "inc.",
		Website:         "",
		Source:          "",
		License:         "",
		Description:     "",
	})
	assert.NoError(err)
	assert.Len(matchIds, 2)
	assert.Contains(matchIds, entry1Created.Id)
	assert.Contains(matchIds, entry2Created.Id)
}
