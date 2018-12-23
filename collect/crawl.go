package collect

import (
	"encoding/json"
	"github.com/michaeldorner/hamster/store"
	"io/ioutil"
	"time"
)

type CrawlRunConfiguration struct {
	ID                string    `json:"id"`
	URL               string    `json:"url"` // remove trailing "/"
	FromDate          time.Time `json:"fromDate"`
	ToDate            time.Time `json:"toDate"`
	OutDir            string    `json:"outDir"`
	Retries           int       `json:"retries"`
	Timeout           int       `json:"timeout"`
	SkipExistingFiles bool      `json:"skipExistingFiles"`
	//HTTPClient        RetryHTTPClient
	//Persistence       store.Persistence
}

var testCrawlRunConfiguration CrawlRunConfiguration = CrawlRunConfiguration{
	ID:                "TESTID",
	URL:               "https://review.openstack.org",
	FromDate:          time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
	ToDate:            time.Date(2016, 1, 2, 0, 0, 0, 0, time.UTC),
	OutDir:            "/TEST/",
	Retries:           10,
	Timeout:           60,
	SkipExistingFiles: true,
}

var TestCrawlRun CrawlRun = NewCrawlRun(testCrawlRunConfiguration)

type CrawlRun struct {
	Config CrawlRunConfiguration
	HTTPClient        RetryHTTPClient
	Persistence       store.Persistence
}

func NewCrawlRun(configuration CrawlRunConfiguration) CrawlRun {
	var crawlRun CrawlRun
	crawlRun.Config = configuration
	crawlRun.Persistence = store.NewPersistence(crawlRun.Config.OutDir, crawlRun.Config.ID)
	crawlRun.HTTPClient = NewRetryHTTPClient(crawlRun.Config.Timeout, crawlRun.Config.Retries, crawlRun.Persistence.LogFile())
	return crawlRun
}


func LoadCrawlRunFile(configurationFilePath string) CrawlRun {
	jsonData, err := ioutil.ReadFile(configurationFilePath)
	if err != nil {
		panic(err)
	}
	var configuration CrawlRunConfiguration
	if err := json.Unmarshal(jsonData, &configuration); err != nil {
		panic(err)
	}
	
	return NewCrawlRun(configuration)
}

func (crawlRun CrawlRun) StoreConfiguration() error {
	data, err := json.MarshalIndent(crawlRun, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(crawlRun.Persistence.CrawlRunFilePath(), data, 0777)
}
