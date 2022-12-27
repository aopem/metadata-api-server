package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"

	"github.com/google/uuid"
)

func ReadFile(filepath string) []byte {
	if !FileExists(filepath) {
		return nil
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil
	}

	return data
}

func WriteFile(filepath string, data []byte) {
	err := os.WriteFile(filepath, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func FileExists(filepath string) bool {
	info, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func GetFolderItems(folderpath string) []fs.DirEntry {
	dir, err := os.Open(folderpath)
	if err != nil {
		return nil
	}

	files, err := dir.ReadDir(0)
	return files
}

func FolderExists(folderpath string) bool {
	info, err := os.Stat(folderpath)

	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func GetUUID() string {
	return uuid.New().String()
}

func CalculateHash(data []byte) string {
	hasher := sha256.New()

	_, err := hasher.Write(data)
	if err != nil {
		// TODO: throw error
		return ""
	}

	return hex.EncodeToString(hasher.Sum(nil))[:8]
}
