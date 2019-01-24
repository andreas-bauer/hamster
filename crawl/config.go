package crawl

import (
	"encoding/json"
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

func UnmarshalConfiguration(jsonData []byte) (Configuration, error) {
	var configuration Configuration
	err := json.Unmarshal(jsonData, &configuration)
	return configuration, err
}
