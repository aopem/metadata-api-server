package services

import (
	"metadata-api-server/internal/brokers"
	"metadata-api-server/internal/testutils"
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataBroker(t *testing.T) {
	testcases := []testutils.Test{{
		Name:     "TestCreateMetadata",
		Function: TestCreateMetadata,
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

	// if folders already exist, clean before running tests
	// then, seed all random numbers that are generated
	utils.DeleteFolder(testutils.TestStorageDirectory)
	utils.DeleteFolder(testutils.TestIndexDirectory)
	testutils.SeedRandomGenerator()
	for i := range testcases {
		t.Run(testcases[i].Name, testcases[i].Function)
	}
}

func TestCreateMetadata(t *testing.T) {
	assert := assert.New(t)

	// create service and dependencies
	mb := brokers.CreateMetadataBroker(testutils.TestStorageDirectory)
	ib := brokers.CreateIndexBroker(testutils.TestIndexDirectory)
	ms := CreateMetadataService(mb, ib)

	// generate data, ID must be empty for new objects
	metadataStore := testutils.GenerateMetadataStore()
	metadataStore.Id = ""

	// create metadata
	metadataStoreCreated, err := ms.CreateMetadata(&metadataStore)
	assert.NoError(err)
	assert.NotNil(metadataStoreCreated)

	// test that it can be retrieved properly
	metadataStoreRetrieved, err := ms.GetMetadataById(metadataStore.Id)
	assert.NoError(err)
	assert.NotNil(metadataStoreRetrieved)
	testutils.AssertMetadataEqual(assert, metadataStoreRetrieved, metadataStoreCreated)

	// cleanup
	ms.DeleteMetadataById(metadataStoreCreated.Id)
}

func TestDeleteMetadataById(t *testing.T) {
	assert := assert.New(t)

	// create service and dependencies
	mb := brokers.CreateMetadataBroker(testutils.TestStorageDirectory)
	ib := brokers.CreateIndexBroker(testutils.TestIndexDirectory)
	ms := CreateMetadataService(mb, ib)

	// generate data, ID must be empty for new objects
	metadataStore := testutils.GenerateMetadataStore()
	metadataStore.Id = ""

	// create metadata
	metadataStoreCreated, err := ms.CreateMetadata(&metadataStore)
	assert.NoError(err)
	assert.NotNil(metadataStoreCreated)

	// test delete, run additional assertions
	_, err = ms.DeleteMetadataById(metadataStoreCreated.Id)
	assert.NoError(err)

	// confirm that a get operation cannot find the metadata
	metadataStoreRetrieved, err := ms.GetMetadataById(metadataStoreCreated.Id)
	assert.Error(err)
	assert.Nil(metadataStoreRetrieved)
}

func TestGetMetadataById(t *testing.T) {
	assert := assert.New(t)

	// create service and dependencies
	mb := brokers.CreateMetadataBroker(testutils.TestStorageDirectory)
	ib := brokers.CreateIndexBroker(testutils.TestIndexDirectory)
	ms := CreateMetadataService(mb, ib)

	// generate data, ID must be empty for new objects
	metadataStore := testutils.GenerateMetadataStore()
	metadataStore.Id = ""

	// create metadata
	metadataStoreCreated, err := ms.CreateMetadata(&metadataStore)
	assert.NoError(err)
	assert.NotNil(metadataStoreCreated)

	// test get by ID
	metadataStoreRetrieved, err := ms.GetMetadataById(metadataStore.Id)
	assert.NoError(err)
	assert.NotNil(metadataStoreRetrieved)
	testutils.AssertMetadataEqual(assert, metadataStoreRetrieved, metadataStoreCreated)

	// cleanup
	ms.DeleteMetadataById(metadataStoreRetrieved.Id)
}

func TestGetMetadata(t *testing.T) {
	assert := assert.New(t)

	// create broker
	mb := brokers.CreateMetadataBroker(testutils.TestStorageDirectory)
	ib := brokers.CreateIndexBroker(testutils.TestIndexDirectory)
	ms := CreateMetadataService(mb, ib)

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

	// retrieve metadata and make sure it matches the initial metadata
	metadataStoreList, err := ms.GetMetadata()
	assert.NoError(err)
	assert.Len(metadataStoreList, len(metadataStoreExpected))
	for i := range metadataStoreList {
		// make sure to compare using ID since not stored in a particular order
		expected := metadataStoreExpected[metadataStoreList[i].Id]
		testutils.AssertMetadataEqual(
			assert,
			&metadataStoreList[i],
			&expected)
	}

	// cleanup
	for i := range metadataStoreList {
		ms.DeleteMetadataById(metadataStoreList[i].Id)
	}
}
