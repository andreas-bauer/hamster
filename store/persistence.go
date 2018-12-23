package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Persistence struct {
	crawlRunID string
	outDir     string
}

func NewPersistence(outDir, crawlRunID string) Persistence {
	var persistence Persistence = Persistence{
		outDir:     outDir,
		crawlRunID: crawlRunID,
	}
	err := os.MkdirAll(filepath.Clean(persistence.outDir+"/"+persistence.crawlRunID+"/data/"), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return persistence
}

func (persistence Persistence) CrawlRunFilePath() string {
	return filepath.Clean(persistence.outDir + "/" + persistence.crawlRunID + "/config.json")
}

func (persistence Persistence) UnitFilePath(id string) string {
	return filepath.Clean(persistence.outDir + "/" + persistence.crawlRunID + fmt.Sprintf("/data/%v.json", id))
}

func (persistence Persistence) UnitFileExists(id string) bool {
	_, err := os.Stat(persistence.UnitFilePath(id))
	return !os.IsNotExist(err)
}

func (persistence Persistence) LogFilePath() string {
	return filepath.Clean(persistence.outDir + "/" + persistence.crawlRunID + "/" + persistence.crawlRunID + ".log")
}

func (persistence Persistence) LogFile() *os.File {
	file, err := os.Create(persistence.LogFilePath())
	if err != nil {
		panic(err)
	}
	return file
}

func (persistence Persistence) StoreUnit(id string, payload []byte) error {
	path := persistence.UnitFilePath(id)
	return ioutil.WriteFile(path, payload, 0644)
}
