package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Repository struct {
	outDir string
}

func NewRepository(outDir string) Repository {
	var repository Repository = Repository{
		outDir: filepath.Clean(outDir),
	}
	err := os.MkdirAll(repository.AppendDataPath(""), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return repository
}

func (repository Repository) OptionsFilePath() string {
	return filepath.Join(repository.outDir, "options.json")
}

func (repository Repository) AppendDataPath(append string) string {
	return filepath.Join(repository.outDir, "data", append)
}

func (repository Repository) FileExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return !os.IsNotExist(err)
}

func (repository Repository) LogFilePath() string {
	return filepath.Join(repository.outDir, "log.log")
}

func (repository Repository) LogFile() *os.File {
	file, err := os.OpenFile(repository.LogFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return file
}

func (repository Repository) Store(path string, payload []byte) error {
	return ioutil.WriteFile(filepath.Clean(path), payload, os.ModePerm)
}
