package crawl

import (
	"github.com/michaeldorner/hamster/pkg/http"
	"github.com/michaeldorner/hamster/pkg/store"
)

type Feed func(Options, http.Client, store.Repository) <-chan Unit
type PostProcess func(Options, http.Client, <-chan Unit) <-chan Unit

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

func filter(options Options, repository store.Repository, in <-chan Unit) <-chan Unit {
	out := make(chan Unit)
	go func() {
		defer close(out)
		for unit := range in {
			if !(options.SkipExistingFiles && repository.UnitFileExists(unit.ID)) {
				out <- unit
			}
		}
	}()
	return out
}

func getPayload(client http.Client, in <-chan Unit) <-chan Unit {
	out := make(chan Unit)
	go func() {
		defer close(out)
		for unit := range in {
			payload, err := client.Get(unit.URL)
			if err != nil {
				panic(err)
			} else {
				unit.Payload = payload
				out <- unit
			}
		}
	}()
	return out
}

func persist(repository store.Repository, in <-chan Unit) {
	for unit := range in {
		path := repository.UnitFilePath(unit.ID)
		err := repository.Store(path, unit.Payload)
		if err != nil {
			panic(err)
		}
	}
}
