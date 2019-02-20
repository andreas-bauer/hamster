package store

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var RepositoryNotEmpty = errors.New("the repository is not empty")

type Repository struct {
	outDir string
}

func NewRepository(outDir string) (*Repository, error) {
	repository := &Repository{
		outDir: filepath.Clean(outDir),
	}
	_, err := os.Stat(filepath.Clean(outDir))
	if !os.IsNotExist(err) {
		return nil, RepositoryNotEmpty
	} else {
		return repository, os.MkdirAll(filepath.Join(repository.outDir, "data"), os.ModePerm)
	}
}

func (repository *Repository) ConfigurationFilePath() string {
	return filepath.Join(repository.outDir, "config.json")
}

func (repository *Repository) LogFile() (*os.File, error) {
	path := filepath.Join(repository.outDir, "log.log")
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
}

func (repository *Repository) StoreItem(fileName string, payload []byte) error {
	path := filepath.Join(repository.outDir, "data", fileName)
	return ioutil.WriteFile(filepath.Clean(path), payload, os.ModePerm)
}

func (repository *Repository) StoreConfiguration(configurationData []byte) error {
	path := filepath.Join(repository.outDir, "config.json")
	return ioutil.WriteFile(filepath.Clean(path), configurationData, os.ModePerm)
}

/*
func FileExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return !os.IsNotExist(err)
}*/
