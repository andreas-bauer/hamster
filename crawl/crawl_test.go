package crawl

import (
	"github.com/michaeldorner/hamster/http"
	"github.com/michaeldorner/hamster/store"
	"io/ioutil"
	"os"
	"path/filepath"
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

var testConfiguration Configuration = Configuration{
	OutDir:           "./repository",
	MaxRetries:       2,
	Timeout:          Duration{time.Duration(10) * time.Second},
	ParallelRequests: 1,
}

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.RemoveAll(testConfiguration.OutDir)
	os.Exit(retCode)
}

func TestCrawl(t *testing.T) {
	Run(testConfiguration, feed)
	_, err := ioutil.ReadFile(filepath.Join(testConfiguration.OutDir, "data", "1.json"))
	if err != nil {
		t.Error(err)
	}
}
