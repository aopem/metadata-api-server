package controllers

import (
	"metadata-api-server/internal/brokers"
	"metadata-api-server/internal/query"
	"metadata-api-server/internal/services"
	"metadata-api-server/internal/testutils"
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestQueryController(t *testing.T) {
	testcases := []testutils.Test{{
		Name:     "TestPutMetadataQuery",
		Function: TestPutMetadataQuery,
	}}

	// if folders already exist, clean before running tests
	// then, seed all random numbers that are generated
	utils.DeleteFolder(testutils.TestStorageDirectory)
	utils.DeleteFolder(testutils.TestIndexDirectory)
	testutils.SeedRandomGenerator()
	for i := range testcases {
		t.Run(testcases[i].Name, testcases[i].Function)
	}
}

func TestPutMetadataQuery(t *testing.T) {
	assert := assert.New(t)

	// create controller and dependencies
	mb := brokers.CreateMetadataBroker(testutils.TestStorageDirectory)
	ib := brokers.CreateIndexBroker(testutils.TestIndexDirectory)
	se := query.CreateSearchEngine(ib)
	qs := services.CreateQueryService(mb, se)
	ms := services.CreateMetadataService(mb, ib)
	qc := CreateQueryController(qs, ms)

	// set gin mode for test, etc. and create a mock router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// set same route as actual server
	router.PUT("/metadata/query", qc.PutMetadataQuery)

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

	// create a query request that should find entry1
	request, err := http.NewRequest(http.MethodPut, "/metadata/query", nil)
	assert.NoError(err)

	// add parameters to request that match entry1
	urlParameters := request.URL.Query()
	urlParameters.Add("title", "Valid App 1")
	urlParameters.Add("company", "Random Inc.")
	urlParameters.Add("website", "https://website.com")
	urlParameters.Add("source", "https://github.com/random/repo")
	request.URL.RawQuery = urlParameters.Encode()

	// execute the request and record the response
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(http.StatusCreated, recorder.Code)

	// verify the response
	response := testutils.QueryResponse{}
	err = yaml.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(err)
	assert.Empty(response.Errors)
	assert.Len(response.Data, 1)
	assert.Equal(response.Data[0], entry1Created.Id)

	// cleanup
	ms.DeleteMetadataById(entry1Created.Id)
	ms.DeleteMetadataById(entry2Created.Id)
}
