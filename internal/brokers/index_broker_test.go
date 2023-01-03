package brokers

import (
	"metadata-api-server/internal/testutils"
	"metadata-api-server/internal/utils"
	"metadata-api-server/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexBroker(t *testing.T) {
	testcases := []testutils.Test{{
		Name:     "TestCreateIndex",
		Function: TestCreateIndex,
	}, {
		Name:     "TestDeleteIndexById",
		Function: TestDeleteIndexById,
	}, {
		Name:     "TestGetIndex",
		Function: TestGetIndex,
	}, {
		Name:     "TestGetIndexPath",
		Function: TestGetIndexPath,
	}, {
		Name:     "TestSaveIndex",
		Function: TestSaveIndex,
	}, {
		Name:     "TestIndexEmpty",
		Function: TestIndexEmpty,
	}, {
		Name:     "TestIndexContains",
		Function: TestIndexContains,
	}}

	// seed all random numbers that are generated
	testutils.SeedRandomGenerator()

	// run tests
	for i := range testcases {
		// if folder already exists, clean before running tests
		utils.DeleteFolder(testutils.TestIndexDirectory)
		t.Run(testcases[i].Name, testcases[i].Function)
	}

	// cleanup
	utils.DeleteFolder(testutils.TestIndexDirectory)
}

func TestCreateIndex(t *testing.T) {
	assert := assert.New(t)

	// create broker, metadata to index
	ib := CreateIndexBroker(testutils.TestIndexDirectory)
	metadataStoreGenerated := testutils.GenerateMetadataStore()

	// index metadata
	ib.CreateIndex(&metadataStoreGenerated)
	assert.True(ib.IndexContains(metadataStoreGenerated.Id))

	// cleanup
	ib.DeleteIndexById(metadataStoreGenerated.Id)
}

func TestDeleteIndexById(t *testing.T) {
	assert := assert.New(t)

	// create broker, metadata to index
	ib := CreateIndexBroker(testutils.TestIndexDirectory)
	metadataStoreGenerated := testutils.GenerateMetadataStore()

	// index metadata
	ib.CreateIndex(&metadataStoreGenerated)
	assert.True(ib.IndexContains(metadataStoreGenerated.Id))

	// assert deletion works properly
	ib.DeleteIndexById(metadataStoreGenerated.Id)
	assert.False(ib.IndexContains(metadataStoreGenerated.Id))
}

func TestGetIndex(t *testing.T) {
	assert := assert.New(t)

	// create broker
	ib := CreateIndexBroker(testutils.TestIndexDirectory)

	// assert that index can be retrieved and is non-nil
	assert.NotNil(ib.GetIndex())
}

func TestGetIndexPath(t *testing.T) {
	assert := assert.New(t)

	// create broker
	ib := CreateIndexBroker(testutils.TestIndexDirectory)

	// assert that index exists at proper location
	assert.NotEmpty(ib.GetIndexPath())
	assert.FileExists(ib.GetIndexPath())
}

func TestSaveIndex(t *testing.T) {
	assert := assert.New(t)

	// create broker, metadata to index
	ib := CreateIndexBroker(testutils.TestIndexDirectory)
	metadataStoreGenerated := testutils.GenerateMetadataStore()

	// index data
	ib.CreateIndex(&metadataStoreGenerated)
	assert.True(ib.IndexContains(metadataStoreGenerated.Id))

	// delete file if it already exists so we know we start with 0 size
	utils.DeleteFile(ib.GetIndexPath())
	assert.True(utils.FileEmpty(ib.GetIndexPath()))

	// now save index
	assert.NoError(ib.SaveIndex())
	bytes, err := utils.ReadFile(ib.GetIndexPath())
	assert.NoError(err)
	assert.NotZero(len(bytes))

	// cleanup
	ib.DeleteIndexById(metadataStoreGenerated.Id)
}

func TestIndexEmpty(t *testing.T) {
	assert := assert.New(t)

	// create broker
	ib := CreateIndexBroker(testutils.TestIndexDirectory)

	// assert that index starts out empty
	assert.True(ib.IndexEmpty())

	// create random metadata and index
	metadataStoreGeneratedList := []models.MetadataStore{}
	for i := 0; i < 10; i++ {
		metadataStoreGenerated := testutils.GenerateMetadataStore()
		metadataStoreGeneratedList = append(metadataStoreGeneratedList, metadataStoreGenerated)
		ib.CreateIndex(&metadataStoreGenerated)
	}

	// assert that new metadata is detected
	assert.False(ib.IndexEmpty())

	// cleanup
	for i := range metadataStoreGeneratedList {
		ib.DeleteIndexById(metadataStoreGeneratedList[i].Id)
	}
}

func TestIndexContains(t *testing.T) {
	assert := assert.New(t)

	// create broker, metadata to index
	ib := CreateIndexBroker(testutils.TestIndexDirectory)
	metadataStoreGenerated := testutils.GenerateMetadataStore()

	// index metadata
	ib.CreateIndex(&metadataStoreGenerated)
	assert.True(ib.IndexContains(metadataStoreGenerated.Id))

	// cleanup
	ib.DeleteIndexById(metadataStoreGenerated.Id)
	assert.False(ib.IndexContains(metadataStoreGenerated.Id))
}
