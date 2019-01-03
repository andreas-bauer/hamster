package crawl

import (
	"encoding/json"
	"io/ioutil"
)

const (
	DefaultTimeout          uint = 120
	DefaultMaxRetryAttempts uint = 5
)

type Options struct {
	URL               string `json:"url"`
	Period            Period `json:"period"`
	OutDir            string `json:"outDir"`
	MaxRetryAttempts  uint   `json:"maxRetryAttempts"`
	Timeout           uint   `json:"timeout"`
	SkipExistingFiles bool   `json:"skipExistingFiles"`
}

func (options Options) JSON() []byte {
	data, err := json.MarshalIndent(options, "", "    ")
	if err != nil {
		panic(err)
	} else {
		return data
	}
}

func LoadOptionsFromJSONFile(configurationFilePath string) Options {
	jsonData, err := ioutil.ReadFile(configurationFilePath)
	if err != nil {
		panic(err)
	}
	var configuration Options
	if err := json.Unmarshal(jsonData, &configuration); err != nil {
		panic(err)
	}

	return configuration
}
