package crawl

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	URL               string `json:"url"`
	Period            Period `json:"period"`
	OutDir            string `json:"outDir"`
	MaxRetryAttempts  uint   `json:"maxRetryAttempts"`
	Timeout           uint   `json:"timeout"`
	SkipExistingFiles bool   `json:"skipExistingFiles"`
	ParallelRequests  uint   `json:"parallelRequests"`
}

func (configuration Configuration) JSON() []byte {
	data, err := json.MarshalIndent(configuration, "", "\t")
	if err != nil {
		panic(err)
	} else {
		return data
	}
}

func LoadConfigurationFromJSONFile(configurationFilePath string) Configuration {
	jsonData, err := ioutil.ReadFile(configurationFilePath)
	if err != nil {
		panic(err)
	}
	var configuration Configuration
	if err := json.Unmarshal(jsonData, &configuration); err != nil {
		panic(err)
	}

	return configuration
}
