package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"log"
	"metadata-api-server/models"
	"os"
	"path/filepath"
	"runtime"
)

func MainDirectory() string {
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		log.Fatal("[ERROR] Could not retrieve main.go directory")
	}

	return filepath.Dir(filepath.Dir(filepath.Dir(filename)))
}

func MetadataEqual(metadataStore *models.MetadataStore, expected *models.MetadataStore) bool {
	equal := true
	equal = equal && metadataStore.Id == expected.Id
	equal = equal && metadataStore.Metadata.Title == expected.Metadata.Title
	equal = equal && metadataStore.Metadata.Version == expected.Metadata.Version
	equal = equal && metadataStore.Metadata.Company == expected.Metadata.Company
	equal = equal && metadataStore.Metadata.Website == expected.Metadata.Website
	equal = equal && metadataStore.Metadata.Source == expected.Metadata.Source
	equal = equal && metadataStore.Metadata.License == expected.Metadata.License
	equal = equal && metadataStore.Metadata.Description == expected.Metadata.Description

	for i := range expected.Metadata.Maintainers {
		equal = equal && metadataStore.Metadata.Maintainers[i].Email == expected.Metadata.Maintainers[i].Email
		equal = equal && metadataStore.Metadata.Maintainers[i].Name == expected.Metadata.Maintainers[i].Name
	}

	return equal
}

func OpenFile(filepath string, flags int, mode fs.FileMode) *os.File {
	file, err := os.OpenFile(filepath, flags, mode)
	if err != nil {
		panic(err)
	}

	return file
}

func CreateFile(filepath string) (*os.File, error) {
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func DeleteFile(filepath string) error {
	return os.Remove(filepath)
}

func ReadFile(filepath string) ([]byte, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func WriteFile(filepath string, data []byte) error {
	return os.WriteFile(filepath, data, os.ModePerm)
}

func FileEmpty(filepath string) bool {
	if !FileExists(filepath) {
		return true
	}

	info, err := os.Stat(filepath)
	if err != nil {
		panic(err)
	}

	return info.Size() == 0
}

func FileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func CreateFolder(folderpath string) {
	if err := os.MkdirAll(folderpath, os.ModePerm); err != nil {
		panic(err)
	}
}

func DeleteFolder(folderpath string) error {
	return os.RemoveAll(folderpath)
}

func FolderExists(folderpath string) bool {
	info, err := os.Stat(folderpath)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func GetFolderItems(folderpath string) ([]fs.DirEntry, error) {
	dir, err := os.Open(folderpath)
	if err != nil {
		return nil, err
	}

	files, err := dir.ReadDir(0)
	return files, err
}

func CalculateHash(data []byte) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil))[:8], nil
}
