package persistence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/michaeldorner/Hamster/collect"
)

type Persistence struct {
	crawlRunID string
	outDir     string
}

func LoadCrawlRunFile(configurationFilePath string) (*collect.CrawlRun, error) {
	var configuration config.CrawlRun
	jsonData, err := ioutil.ReadFile(configurationFilePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonData, &configuration); err != nil {
		return nil, err
	}
	return &configuration, nil
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

func (persistence Persistence) StoreUnit(id string, payload []byte) error {
	path := persistence.UnitFilePath(id)
	return ioutil.WriteFile(path, payload, 0644)
}

func (persistence Persistence) StoreCrawlRun(configuration config.CrawlRun) error {
	data, err := json.MarshalIndent(configuration, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(persistence.CrawlRunFilePath(), data, os.ModePerm)
}
