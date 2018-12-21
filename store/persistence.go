package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Persistence struct {
	crawlRunID string
	outDir     string
}

func LoadConfigurationFile(configurationFilePath string) (*Configuration, error) {
	var configuration Configuration
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

func (persistence Persistence) ConfigurationFilePath() string {
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

func (persistence Persistence) StoreConfiguration(configuration Configuration) error {
	data, err := json.MarshalIndent(configuration, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(persistence.ConfigurationFilePath(), data, os.ModePerm)
}
