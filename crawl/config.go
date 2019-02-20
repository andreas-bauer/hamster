package crawl

import (
	"encoding/json"
)

type Configuration struct {
	URL               string   `json:"url"`
	Period            Period   `json:"period"`
	OutDir            string   `json:"outDir"`
	MaxRetries        uint     `json:"maxRetries"`
	Timeout           Duration `json:"timeout"`
	SkipExistingFiles bool     `json:"skipExistingFiles"`
	ParallelRequests  uint     `json:"parallelRequests"`
}

func (configuration Configuration) JSON() ([]byte, error) {
	return json.MarshalIndent(configuration, "", "\t")
}

func UnmarshalConfiguration(jsonData []byte) (Configuration, error) {
	var configuration Configuration
	err := json.Unmarshal(jsonData, &configuration)
	return configuration, err
}
