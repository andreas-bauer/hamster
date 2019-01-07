package crawl

import (
	"github.com/michaeldorner/hamster/pkg/http"
	"github.com/michaeldorner/hamster/pkg/store"
)

type Feed func(Options, http.Client, store.Repository) <-chan Item
type PostProcess func(Options, http.Client, <-chan Item) <-chan Item

func Run(options Options, feed Feed, postProcess PostProcess) {
	repository := store.NewRepository(options.OutDir)
	client := http.NewClient(options.Timeout, options.MaxRetryAttempts, repository.LogFile())

	storeOptions(options, repository)
	afterFeed := feed(options, client, repository)
	afterFilter := filter(options, repository, afterFeed)
	afterPayload := getPayload(client, afterFilter)
	afterPostProcess := postProcess(options, client, afterPayload)
	persist(repository, afterPostProcess)
}

func storeOptions(options Options, repository store.Repository) {
	jsonData := options.JSON()
	path := repository.OptionsFilePath()
	err := repository.Store(path, jsonData)
	if err != nil {
		panic(err)
	}
}

func filter(options Options, repository store.Repository, in <-chan Item) <-chan Item {
	out := make(chan Item)
	go func() {
		defer close(out)
		for item := range in {
			path := repository.AppendDataPath(item.FileName()) 
			if !(options.SkipExistingFiles && repository.FileExists(path)) {
				out <- item
			}
		}
	}()
	return out
}

func getPayload(client http.Client, in <-chan Item) <-chan Item {
	out := make(chan Item)
	go func() {
		defer close(out)
		for item := range in {
			payload, err := client.Get(item.URL)
			if err != nil {
				panic(err)
			} else {
				item.Payload = payload
				out <- item
			}
		}
	}()
	return out
}

func persist(repository store.Repository, in <-chan Item) {
	for item := range in {
		file_path := repository.AppendDataPath(item.FileName())
		err := repository.Store(file_path, item.Payload)
		if err != nil {
			panic(err)
		}
	}
}
