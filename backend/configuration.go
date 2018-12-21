package internal

import (
	"errors"
	"time"
)

var ToolNotSupported error = errors.New("Gerry supports only Gerry and Phabricator")

type ReviewTool string

const (
	Gerrit      ReviewTool = "Gerrit"
	Phabricator ReviewTool = "Phabricator"
)

type Configuration struct {
	CrawlRunID        string     `json:"crawlRunID"`
	URL               string     `json:"url"` // remove trailing "/"
	FromDate          time.Time  `json:"fromDate"`
	ToDate            time.Time  `json:"toDate"`
	OutDir            string     `json:"outDir"`
	Retries           int        `json:"retries"`
	TimeOut           int        `json:"timeOut"`
	SkipExistingFiles bool       `json:"skipExistingFiles"`
	ReviewTool        ReviewTool `json:"reviewTool"`
}

var TestConfiguration Configuration = Configuration{
	CrawlRunID:        "TEST",
	URL:               "https://review.openstack.org",
	FromDate:          time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
	ToDate:            time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
	OutDir:            "/Users/michaeldorner/Desktop/Gerry/",
	Retries:           2,
	TimeOut:           60,
	SkipExistingFiles: true,
	ReviewTool:        "Gerrit",
}
