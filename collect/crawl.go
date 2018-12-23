package collect

import (
	"encoding/json"
	"github.com/michaeldorner/hamster/store"
	"io/ioutil"
	"time"
)

type CrawlRun struct {
	ID                string    `json:"id"`
	URL               string    `json:"url"` // remove trailing "/"
	FromDate          time.Time `json:"fromDate"`
	ToDate            time.Time `json:"toDate"`
	OutDir            string    `json:"outDir"`
	Retries           int       `json:"retries"`
	Timeout           int       `json:"timeout"`
	SkipExistingFiles bool      `json:"skipExistingFiles"`
	HTTPClient        RetryHTTPClient
	Persistence       store.Persistence
}

var standardPersistence store.Persistence = store.NewPersistence("/Users/michaeldorner/Desktop/Gerry/", "TEST")
var standardHTTPClient = NewRetryHTTPClient(60, 2, standardPersistence.LogFile())

var TestCrawlRun CrawlRun = CrawlRun{
	ID:                "TEST",
	URL:               "https://review.openstack.org",
	FromDate:          time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
	ToDate:            time.Date(2016, 1, 2, 0, 0, 0, 0, time.UTC),
	OutDir:            "/Users/michaeldorner/Desktop/Gerry/",
	Retries:           10,
	Timeout:           60,
	SkipExistingFiles: true,
	HTTPClient:        standardHTTPClient,
	Persistence:       standardPersistence,
}

func LoadCrawlRunFile(configurationFilePath string) CrawlRun {
	jsonData, err := ioutil.ReadFile(configurationFilePath)
	if err != nil {
		panic(err)
	}
	var crawlRun CrawlRun
	if err := json.Unmarshal(jsonData, &crawlRun); err != nil {
		panic(err)
	}
	crawlRun.Persistence = store.NewPersistence(crawlRun.OutDir, crawlRun.ID)
	crawlRun.HTTPClient = NewRetryHTTPClient(crawlRun.Timeout, crawlRun.Retries, crawlRun.Persistence.LogFile())

	return crawlRun
}

func (crawlRun CrawlRun) StoreConfiguration() error {
	data, err := json.MarshalIndent(crawlRun, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(crawlRun.Persistence.CrawlRunFilePath(), data, 0777)
}
