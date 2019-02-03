package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Repository struct {
	outDir string
}

func NewRepository(outDir string) (Repository, error) {
	var repository Repository = Repository{
		outDir: filepath.Clean(outDir),
	}
	return repository, os.MkdirAll(repository.AppendDataPath(""), os.ModePerm)
}

func (repository Repository) ConfigurationFilePath() string {
	return filepath.Join(repository.outDir, "config.json")
}

func (repository Repository) AppendDataPath(append string) string {
	return filepath.Join(repository.outDir, "data", append)
}

func (repository Repository) FileExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return !os.IsNotExist(err)
}

func (repository Repository) logFilePath() string {
	return filepath.Join(repository.outDir, "crawl.log")
}

func (repository Repository) LogFile() (*os.File, error) {
	return os.OpenFile(repository.logFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
}

func (repository Repository) Store(path string, payload []byte) error {
	return ioutil.WriteFile(filepath.Clean(path), payload, os.ModePerm)
}
