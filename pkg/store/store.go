package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Repository struct {
	outDir string
}

func NewRepository(outDir string) Repository {
	var repository Repository = Repository{
		outDir: outDir,
	}
	err := os.MkdirAll(filepath.Clean(repository.outDir+"/data/"), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return repository
}

func (repository Repository) OptionsFilePath() string {
	return filepath.Clean(repository.outDir + "/options.json")
}

func (repository Repository) DataPath() string {
	return filepath.Clean(repository.outDir + fmt.Sprintf("/data/"))
}

func (repository Repository) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (repository Repository) LogFilePath() string {
	return filepath.Clean(repository.outDir + "/log.log")
}

func (repository Repository) LogFile() *os.File {
	file, err := os.OpenFile(repository.LogFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return file
}

func (repository Repository) Store(path string, payload []byte) error {
	return ioutil.WriteFile(path, payload, os.ModePerm)
}
