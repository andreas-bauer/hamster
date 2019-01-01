package crawl

import (
	"fmt"
	"testing"
)


var TestOptions Options = Options{
	URL:               "https://review.openstack.org",
	FromDate:          NewDate(2018, 1, 1),
	ToDate:            NewDate(2018, 1, 31),
	OutDir:            "./",
	MaxRetryAttempts:  10,
	Timeout:           120,
	SkipExistingFiles: false,
}

func TestTs(t *testing.T) {
	jsonString := TestOptions.JSON()
	fmt.Println(jsonString)
}
