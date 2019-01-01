package crawl

import (
	"encoding/json"
	"io/ioutil"
)

type Options struct {
	URL               string    `json:"url"` // remove trailing "/"
	FromDate          Date `json:"fromDate"`
	ToDate            Date `json:"toDate"`
	OutDir            string    `json:"outDir"`
	MaxRetryAttempts  int       `json:"maxRetryAttempts"`
	Timeout           int       `json:"timeout"`
	SkipExistingFiles bool      `json:"skipExistingFiles"`
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