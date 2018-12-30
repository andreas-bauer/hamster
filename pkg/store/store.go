package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Repository struct {
	crawlTaskID string
	outDir      string
}

func NewRepository(outDir, crawlTaskID string) Repository {
	var repository Repository = Repository{
		outDir:      outDir,
		crawlTaskID: crawlTaskID,
	}
	err := os.MkdirAll(filepath.Clean(repository.outDir+"/"+repository.crawlTaskID+"/data/"), 0644)
	if err != nil {
		panic(err)
	}
	return repository
}

func (repository Repository) OptionsFilePath() string {
	return filepath.Clean(repository.outDir + "/" + repository.crawlTaskID + "/options.json")
}

func (repository Repository) UnitFilePath(id string) string {
	return filepath.Clean(repository.outDir + "/" + repository.crawlTaskID + fmt.Sprintf("/data/%v.json", id))
}

func (repository Repository) UnitFileExists(id string) bool {
	_, err := os.Stat(repository.UnitFilePath(id))
	return !os.IsNotExist(err)
}

func (repository Repository) LogFilePath() string {
	return filepath.Clean(repository.outDir + "/" + repository.crawlTaskID + "/" + repository.crawlTaskID + ".log")
}

func (repository Repository) LogFile() *os.File {
	file, err := os.Create(repository.LogFilePath())
	if err != nil {
		panic(err)
	}
	return file
}

func (repository Repository) Store(path string, payload []byte) error {
	return ioutil.WriteFile(path, payload, 0644)
}
