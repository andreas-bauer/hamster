package crawl

import (
	"fmt"
	"github.com/michaeldorner/hamster/http"
	"github.com/michaeldorner/hamster/store"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var testConfiguration Configuration = Configuration{
	URL:              "https://jsonplaceholder.typicode.com/todos/",
	OutDir:           "./repository",
	MaxRetries:       2,
	Timeout:          Duration{time.Duration(10) * time.Second},
	ParallelRequests: 1,
}

var feed Feed = func(configuration Configuration, client http.Client, repository *store.Repository) <-chan Item {
	items := make(chan Item)
	go func() {
		defer close(items)
		for i := 1; i < 3; i++ {
			id := fmt.Sprint(i)
			items <- Item{
				ID:                 id,
				URL:                configuration.URL + id,
				FileNameExtensions: "json",
			}
		}

	}()
	return items
}

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.RemoveAll(testConfiguration.OutDir)
	os.Exit(retCode)
}

func TestCrawl(t *testing.T) {
	Run(testConfiguration, feed)
	for i := 1; i < 3; i++ {
		_, err := ioutil.ReadFile(filepath.Join(testConfiguration.OutDir, "data", fmt.Sprintf("%v.json", i)))
		if err != nil {
			t.Error(err)
		}
	}
}
