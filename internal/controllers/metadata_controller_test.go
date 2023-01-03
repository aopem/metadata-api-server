package controllers

import (
	"bytes"
	"metadata-api-server/internal/brokers"
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

func TestMetadataController(t *testing.T) {
	testcases := []testutils.Test{{
		Name:     "TestPutMetadata",
		Function: TestPutMetadata,
	}, {
		Name:     "TestDeleteMetadataById",
		Function: TestDeleteMetadataById,
	}, {
		Name:     "TestGetMetadataById",
		Function: TestGetMetadataById,
	}, {
		Name:     "TestGetMetadata",
		Function: TestGetMetadata,
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

func TestPutMetadata(t *testing.T) {
	assert := assert.New(t)

	// create controller and dependencies
	mb := brokers.CreateMetadataBroker(testutils.TestStorageDirectory)
	ib := brokers.CreateIndexBroker(testutils.TestIndexDirectory)
	ms := services.CreateMetadataService(mb, ib)
	mc := CreateMetadataController(ms)

	// set gin mode for test, etc. and create a mock router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// set same route as actual server
	router.PUT("/metadata", mc.PutMetadata)

	// create test metadata to send as body
	metadataStore := &models.MetadataStore{
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

	// put into a buffer for the request body
	body := bytes.Buffer{}
	err := yaml.NewEncoder(&body).Encode(metadataStore)
	assert.NoError(err)

	// create mock request
	request, err := http.NewRequest(http.MethodPut, "/metadata", &body)
	assert.NoError(err)

	// execute the request and record the response
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(http.StatusCreated, recorder.Code)

	// verify the response
	response := testutils.MetadataResponse{}
	err = yaml.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(err)
	assert.Empty(response.Errors)
	assert.Len(response.Data, 1)
	testutils.AssertMetadataEqual(assert, &response.Data[0], metadataStore)

	// cleanup
	ms.DeleteMetadataById(response.Data[0].Id)
}

func TestDeleteMetadataById(t *testing.T) {
	assert := assert.New(t)

	// create controller and dependencies
	mb := brokers.CreateMetadataBroker(testutils.TestStorageDirectory)
	ib := brokers.CreateIndexBroker(testutils.TestIndexDirectory)
	ms := services.CreateMetadataService(mb, ib)
	mc := CreateMetadataController(ms)

	// set gin mode for test, etc. and create a mock router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// set same route as actual server
	router.DELETE("/metadata/:id", mc.DeleteMetadataById)

	// create test metadata that will be deleted
	metadataStore := &models.MetadataStore{
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
	metadataStoreCreated, err := ms.CreateMetadata(metadataStore)
	assert.NoError(err)

	// create mock request
	request, err := http.NewRequest(http.MethodDelete, "/metadata/"+metadataStoreCreated.Id, nil)
	assert.NoError(err)

	// execute the request and record the response
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(http.StatusGone, recorder.Code)

	// verify the response
	response := testutils.MetadataResponse{}
	err = yaml.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(err)
	assert.Empty(response.Errors)

	// confirm that the record cannot be retrieved
	metadataStoreRetrieved, err := ms.GetMetadataById(metadataStoreCreated.Id)
	assert.Error(err)
	assert.Nil(metadataStoreRetrieved)
}

func TestGetMetadataById(t *testing.T) {
	assert := assert.New(t)

	// create controller and dependencies
	mb := brokers.CreateMetadataBroker(testutils.TestStorageDirectory)
	ib := brokers.CreateIndexBroker(testutils.TestIndexDirectory)
	ms := services.CreateMetadataService(mb, ib)
	mc := CreateMetadataController(ms)

	// set gin mode for test, etc. and create a mock router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// set same route as actual server
	router.GET("/metadata/:id", mc.GetMetadataById)

	// create test metadata that will be retrieved
	metadataStore := &models.MetadataStore{
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
	metadataStoreCreated, err := ms.CreateMetadata(metadataStore)
	assert.NoError(err)

	// create mock request
	request, err := http.NewRequest(http.MethodGet, "/metadata/"+metadataStoreCreated.Id, nil)
	assert.NoError(err)

	// execute the request and record the response
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(http.StatusOK, recorder.Code)

	// verify the response
	response := testutils.MetadataResponse{}
	err = yaml.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(err)
	assert.Empty(response.Errors)
	assert.Len(response.Data, 1)
	testutils.AssertMetadataAndIdEqual(assert, &response.Data[0], metadataStoreCreated)

	// cleanup
	ms.DeleteMetadataById(response.Data[0].Id)
}

func TestGetMetadata(t *testing.T) {
	assert := assert.New(t)

	// create controller and dependencies
	mb := brokers.CreateMetadataBroker(testutils.TestStorageDirectory)
	ib := brokers.CreateIndexBroker(testutils.TestIndexDirectory)
	ms := services.CreateMetadataService(mb, ib)
	mc := CreateMetadataController(ms)

	// set gin mode for test, etc. and create a mock router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// set same route as actual server
	router.GET("/metadata", mc.GetMetadata)

	// create test metadata that will be retrieved
	// generate/save a list of metadata
	listLength := 10
	metadataStoreExpected := map[string]models.MetadataStore{}
	for i := 0; i < listLength; i++ {
		// save generated metadata to a map, where key is ID
		metadataStoreGenerated := testutils.GenerateMetadataStore()
		metadataStoreGenerated.Id = ""
		metadataStoreCreated, err := ms.CreateMetadata(&metadataStoreGenerated)
		metadataStoreExpected[metadataStoreCreated.Id] = *metadataStoreCreated
		assert.NoError(err)
		assert.NotNil(metadataStoreCreated)
	}

	// create mock request
	request, err := http.NewRequest(http.MethodGet, "/metadata", nil)
	assert.NoError(err)

	// execute the request and record the response
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(http.StatusOK, recorder.Code)

	// verify the response
	response := testutils.MetadataResponse{}
	err = yaml.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(err)
	assert.Empty(response.Errors)
	assert.Len(response.Data, listLength)
	for i := range response.Data {
		expected := metadataStoreExpected[response.Data[i].Id]
		testutils.AssertMetadataAndIdEqual(assert, &response.Data[i], &expected)
	}

	// cleanup
	for i := range response.Data {
		ms.DeleteMetadataById(response.Data[i].Id)
	}
}
