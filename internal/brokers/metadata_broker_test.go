package brokers

import (
	"metadata-api-server/internal/testutils"
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testStorageDirectory = filepath.Join(utils.MainDirectory(), "localStoreTest")

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
		Name:     "TestGetMetadataList",
		Function: TestGetMetadataList,
	}}

	// if folder already exists, clean before running tests
	// then, seed all random numbers that are generated
	utils.DeleteFolder(testStorageDirectory)
	testutils.SeedRandomGenerator()
	for i := range testcases {
		t.Run(testcases[i].Name, testcases[i].Function)
	}
}

func TestCreateMetadata(t *testing.T) {
	assert := assert.New(t)

	// create data, broker
	mb := CreateMetadataBroker(testStorageDirectory)
	metadataStore := testutils.GenerateMetadataStore()

	// test function, run assertions
	metadataStoreCreated, err := mb.CreateMetadata(&metadataStore)
	metadataFile := filepath.Join(testStorageDirectory, metadataStoreCreated.Id+".yaml")
	assert.NoError(err)
	assert.FileExists(metadataFile)
	utils.DeleteFile(metadataFile)
}

func TestDeleteMetadataById(t *testing.T) {
	assert := assert.New(t)

	// create data, broker
	mb := CreateMetadataBroker(testStorageDirectory)
	metadataStore := testutils.GenerateMetadataStore()
	_, err := mb.CreateMetadata(&metadataStore)
	metadataFile := filepath.Join(testStorageDirectory, metadataStore.Id+".yaml")
	assert.NoError(err)
	assert.FileExists(metadataFile)

	// test delete, run additional assertions
	_, err = mb.DeleteMetadataById(metadataStore.Id)
	assert.NoError(err)
	assert.NoFileExists(metadataFile)
}

func TestGetMetadataById(t *testing.T) {
	assert := assert.New(t)

	// create data, broker
	mb := CreateMetadataBroker(testStorageDirectory)
	metadataStore := testutils.GenerateMetadataStore()
	_, err := mb.CreateMetadata(&metadataStore)
	metadataFile := filepath.Join(testStorageDirectory, metadataStore.Id+".yaml")
	assert.NoError(err)
	assert.FileExists(metadataFile)

	// test get by ID and assert
	metadataStoreCreated, err := mb.GetMetadataById(metadataStore.Id)
	assert.NoError(err)
	testutils.AssertMetadataEqual(assert, metadataStoreCreated, &metadataStore)

	// cleanup
	utils.DeleteFile(metadataFile)
}

func TestGetMetadataList(t *testing.T) {
	assert := assert.New(t)

	// create broker
	mb := CreateMetadataBroker(testStorageDirectory)

	// generate/save a list of metadata
	listLength := 10
	metadataStoreExpected := map[string]models.MetadataStore{}
	for i := 0; i < listLength; i++ {
		// save generated metadata to a map, where key is ID
		metadataStoreGenerated := testutils.GenerateMetadataStore()
		metadataStoreExpected[metadataStoreGenerated.Id] = metadataStoreGenerated
		_, err := mb.CreateMetadata(&metadataStoreGenerated)
		assert.NoError(err)
	}

	// retrieve metadata and make sure it matches the initial metadata
	metadataStoreList, err := mb.GetMetadataList()
	assert.NoError(err)
	for i := range metadataStoreList {
		// make sure to compare using ID since not stored in a particular order
		expected := metadataStoreExpected[metadataStoreList[i].Id]
		testutils.AssertMetadataEqual(
			assert,
			&metadataStoreList[i],
			&expected)
	}

	// cleanup
	utils.DeleteFolder(testStorageDirectory)
}
