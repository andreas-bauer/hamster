package crawl

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Options struct {
	ID                string    `json:"id"`
	URL               string    `json:"url"` // remove trailing "/"
	FromDate          time.Time `json:"fromDate"`
	ToDate            time.Time `json:"toDate"`
	OutDir            string    `json:"outDir"`
	Retries           int       `json:"retries"`
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

var TestOptions Options = Options{
	ID:                "TESTID",
	URL:               "https://review.openstack.org",
	FromDate:          time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
	ToDate:            time.Date(2016, 1, 2, 0, 0, 0, 0, time.UTC),
	OutDir:            "/Users/michaeldorner/Desktop/",
	Retries:           10,
	Timeout:           60,
	SkipExistingFiles: true,
}

func LoadTaskFile(configurationFilePath string) Options {
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

/*
func (options Options) RepositoryConfiguration() error {
	data := options.JSON()
	return crawlTask.Repository.RepositoryConfiguration(data)
}
*/
