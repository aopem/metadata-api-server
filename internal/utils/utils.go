package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
)

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
