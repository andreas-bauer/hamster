package crawl

import (
	"encoding/json"
	"time"
)

type jsonConfiguration struct {
	URL              string      `json:"url"`
	Feed             interface{} `json:"feed"`
	OutDir           string      `json:"outDir"`
	MaxRetries       uint        `json:"maxRetries"`
	Timeout          Duration    `json:"timeout"`
	ParallelRequests uint        `json:"parallelRequests"`
}

type Configuration struct {
	jsonConfiguration
}

func (config *Configuration) URL() string {
	return config.jsonConfiguration.URL
}

func (config *Configuration) Feed() interface{} {
	return config.jsonConfiguration.Feed
}

func (config *Configuration) OutDir() string {
	return config.jsonConfiguration.OutDir
}

func (config *Configuration) MaxRetries() uint {
	return config.jsonConfiguration.MaxRetries
}

func (config *Configuration) Timeout() time.Duration {
	return config.jsonConfiguration.Timeout.Duration
}

func (config *Configuration) ParallelRequests() uint {
	return config.jsonConfiguration.ParallelRequests
}

func (config *Configuration) JSON() ([]byte, error) {
	return json.MarshalIndent(config.jsonConfiguration, "", "\t")
}
func (config *Configuration) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &config.jsonConfiguration)
}

func UnmarshalConfiguration(jsonData []byte) (*Configuration, error) {
	var config Configuration
	err := json.Unmarshal(jsonData, &config)
	return &config, err
}
