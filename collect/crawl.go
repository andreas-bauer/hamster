package collect

import (
	"time"
	"github.com/michaeldorner/hamster/persistence"
)


type CrawlRun struct {
	CrawlRunID        string     `json:"crawlRunID"`
	URL               string     `json:"url"` // remove trailing "/"
	FromDate          time.Time  `json:"fromDate"`
	ToDate            time.Time  `json:"toDate"`
	OutDir            string     `json:"outDir"`
	Retries           int        `json:"retries"`
	TimeOut           int        `json:"timeOut"`
	SkipExistingFiles bool       `json:"skipExistingFiles"`
	Client 				RetryHTTPClient
	persistence			store.Persistence
} 

var TestCrawlRun CrawlRun = CrawlRun{
	CrawlRunID:        "TEST",
	URL:               "https://review.openstack.org",
	FromDate:          time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
	ToDate:            time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
	OutDir:            "/Users/michaeldorner/Desktop/Gerry/",
	Retries:           2,
	TimeOut:           60,
	SkipExistingFiles: true,
}

func CrawlFromGerry(filePath string)

	crawlRun, err := persistence.LoadCrawlRunFile(filePath)
	if err != nil {
		panic(err)
	}

	channel_1 := gerrit.Generate(crawlRun)
	channel_2 := preset.Filter(channel_1, crawlRun)
	channel_3 := preset.GetPayload(channel_2, crawlRun)
	channel_4 := gerrit.PostProcess(channel_3, crawlRun)
	
	preset.store(channel_4, crawlRun)
}


func CrawlFromPhabricator(filePath string)

	crawlRun, err := persistence.LoadCrawlRunFile(filePath)
	if err != nil {
		panic(err)
	}
	channel_1 := phabricator.Generate(crawlRun)
	channel_2 := preset.Filter(channel_1, crawlRun)
	channel_3 := preset.GetPayload(channel_2, crawlRun)
	channel_4 := phabricator.PostProcess(channel_3, crawlRun)
	
	preset.store(channel_4, crawlRun)
}