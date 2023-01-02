package testutils

import (
	"math/rand"
	"metadata-api-server/models"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Name     string
	Function func(t *testing.T)
}

func AssertMetadataEqual(assert *assert.Assertions, metadataStore *models.MetadataStore, expected *models.MetadataStore) {
	assert.Equal(metadataStore.Id, expected.Id)
	assert.Equal(metadataStore.Metadata.Title, expected.Metadata.Title)
	assert.Equal(metadataStore.Metadata.Version, expected.Metadata.Version)
	assert.Equal(metadataStore.Metadata.Company, expected.Metadata.Company)
	assert.Equal(metadataStore.Metadata.Website, expected.Metadata.Website)
	assert.Equal(metadataStore.Metadata.Source, expected.Metadata.Source)
	assert.Equal(metadataStore.Metadata.License, expected.Metadata.License)
	assert.Equal(metadataStore.Metadata.Description, expected.Metadata.Description)

	for i := range expected.Metadata.Maintainers {
		assert.Equal(metadataStore.Metadata.Maintainers[i].Email, expected.Metadata.Maintainers[i].Email)
		assert.Equal(metadataStore.Metadata.Maintainers[i].Name, expected.Metadata.Maintainers[i].Name)
	}
}

func GenerateMetadataStore() models.MetadataStore {
	metadata := models.Metadata{
		Title:       generateRandomString(generateRandomNumber(0, 100)),
		Version:     generateRandomString(generateRandomNumber(0, 100)),
		Company:     generateRandomString(generateRandomNumber(0, 100)),
		Website:     generateRandomString(generateRandomNumber(0, 100)),
		Source:      generateRandomString(generateRandomNumber(0, 100)),
		License:     generateRandomString(generateRandomNumber(0, 100)),
		Description: generateRandomString(generateRandomNumber(0, 100)),
		Maintainers: []models.Maintainer{},
	}

	for i := 0; i < generateRandomNumber(0, 10); i++ {
		maintainer := models.Maintainer{
			Name:  generateRandomString(generateRandomNumber(0, 100)),
			Email: generateRandomString(generateRandomNumber(0, 100)),
		}
		metadata.Maintainers = append(metadata.Maintainers, maintainer)
	}

	return models.MetadataStore{
		Id:       generateUUID(),
		Metadata: &metadata,
	}
}

func SeedRandomGenerator() {
	rand.Seed(time.Now().UnixNano())
}

func generateUUID() string {
	return uuid.New().String()[:8]
}

func generateRandomNumber(low int, high int) int {
	return low + rand.Intn(high-low)
}

func generateRandomString(length int) string {
	// declare letters that will be used
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	// pick a random letter of the alphabet for each character of the string
	var sb strings.Builder
	for i := 0; i < length; i++ {
		ch := alphabet[rand.Intn(len(alphabet))]
		sb.WriteRune(ch)
	}

	s := sb.String()
	return s
}
