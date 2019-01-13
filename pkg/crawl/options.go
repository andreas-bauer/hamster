package crawl

import (
	"encoding/json"
	"io/ioutil"
)

const (
	DefaultTimeout          uint = 120
	DefaultMaxRetryAttempts uint = 5
)

type Configuration struct {
	URL               string `json:"url,omitempty"`
	Period            Period `json:"period,omitempty"`
	OutDir            string `json:"outDir,omitempty"`
	MaxRetryAttempts  uint   `json:"maxRetryAttempts,omitempty"`
	Timeout           uint   `json:"timeout,omitempty"`
	SkipExistingFiles bool   `json:"skipExistingFiles,omitempty"`
}

func (options Configuration) JSON() []byte {
	data, err := json.MarshalIndent(options, "", "    ")
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
