package crawl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/michaeldorner/hamster/http"
	"github.com/michaeldorner/hamster/store"
)

var feed Feed = func(configuration *Configuration, client *http.Client, repository *store.Repository) <-chan Item {
	items := make(chan Item)
	go func() {
		defer close(items)
		for i := 1; i < 3; i++ {
			id := fmt.Sprint(i)
			req, err := http.NewGetRequest(configuration.URL() + id)
			if err != nil {
				panic(err)
			}
			items <- Item{
				ID:                 id,
				Request:            req,
				FileNameExtensions: "json",
			}
		}

	}()
	return items
}

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.RemoveAll("./repository/")
	os.Exit(retCode)
}

func TestCrawl(t *testing.T) {
	var feed Feed = func(configuration *Configuration, client *http.Client, repository *store.Repository) <-chan Item {
		items := make(chan Item)
		go func() {
			defer close(items)
			for i := 1; i < 3; i++ {
				id := fmt.Sprint(i)
				req, err := http.NewGetRequest(configuration.URL() + id)
				if err != nil {
					t.Error(err)
				}
				items <- Item{
					ID:                 id,
					Request:            req,
					FileNameExtensions: "json",
				}
			}

		}()
		return items
	}

	var configurationJSONData = `{
		"url": "https://jsonplaceholder.typicode.com/todos/",
		"feed": {},
		"outDir": "./repository/",
		"maxRetries": 2,
		"timeout": "10s",
		"parallelRequests": 1
	}`

	testConfiguration, err := UnmarshalConfiguration([]byte(configurationJSONData))
	if err != nil {
		t.Error(err)
	}

	Run(testConfiguration, feed)
	for i := 1; i < 3; i++ {
		_, err := ioutil.ReadFile(filepath.Join(testConfiguration.OutDir(), "data", fmt.Sprintf("%v.json", i)))
		if err != nil {
			t.Error(err)
		}
	}
}
