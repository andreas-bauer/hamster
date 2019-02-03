package crawl

import (
	"github.com/michaeldorner/hamster/http"
	"github.com/michaeldorner/hamster/store"
	"testing"
)

var postProcess PostProcess = func(configuration Configuration, client http.Client, in <-chan Item) <-chan Item {
	items := make(chan Item)
	go func() {
		defer close(items)
		for item := range in {
			items <- item
		}
	}()
	return items
}

var feed Feed = func(Configuration, http.Client, store.Repository) <-chan Item {
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
		OutDir:            "./",
		MaxRetryAttempts:  2,
		Timeout:           10,
		SkipExistingFiles: false,
		ParallelRequests:  1,
	}
	Run(configuration, feed, postProcess)
}
