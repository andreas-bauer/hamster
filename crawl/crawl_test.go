package crawl

import (
	"github.com/michaeldorner/hamster/http"
	"github.com/michaeldorner/hamster/store"
	"os"
	"testing"
	"time"
)

var feed Feed = func(Configuration, http.Client, *store.Repository) <-chan Item {
	items := make(chan Item)
	go func() {
		defer close(items)
		items <- Item{
			ID:                 "1",
			URL:                "https://jsonplaceholder.typicode.com/todos/1",
			FileNameExtensions: "json",
		}
	}()
	return items
}

func TestCrawl(t *testing.T) {
	configuration := Configuration{
		OutDir:            "./repository",
		MaxRetries:        2,
		Timeout:           Duration{time.Duration(10) * time.Second},
		SkipExistingFiles: false,
		ParallelRequests:  1,
	}
	Run(configuration, feed)
	//check for file 1.json
	os.RemoveAll("./repository")
}
