package crawl

import (
	"os"
	"sync"

	"github.com/michaeldorner/hamster/http"
	"github.com/michaeldorner/hamster/store"
)

type Feed func(Configuration, http.Client, store.Repository) <-chan Item
type PostProcess func(Configuration, http.Client, <-chan Item) <-chan Item

func Run(options Configuration, feed Feed, postProcess PostProcess) {
	repository := store.NewRepository(options.OutDir)

	client := http.NewClient(options.Timeout, options.MaxRetryAttempts)

	storeConfiguration(options, repository)
	afterFeed := feed(options, client, repository)
	afterFilter := filter(options, repository, afterFeed)
	afterPayload := getPayload(client, afterFilter, options.ParallelRequests, repository.LogFile())
	afterPostProcess := postProcess(options, client, afterPayload)
	persist(repository, afterPostProcess)
}

func storeConfiguration(options Configuration, repository store.Repository) {
	jsonData := options.JSON()
	path := repository.ConfigurationFilePath()
	err := repository.Store(path, jsonData)
	if err != nil {
		panic(err)
	}
}

func filter(options Configuration, repository store.Repository, in <-chan Item) <-chan Item {
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

func getPayload(client http.Client, in <-chan Item, numParallelRequests uint, logFile *os.File) <-chan Item {
	out := make(chan Item)

	go func() {
		log := make(chan string)

		var dowloadWaitGroup sync.WaitGroup
		var logWaitGroup sync.WaitGroup

		for i := uint(0); i < numParallelRequests; i++ {
			dowloadWaitGroup.Add(1)
			go func() {
				defer dowloadWaitGroup.Done()
				for item := range in {
					payload, err := client.Get(item.URL, log)
					if err != nil {
						panic(err)
					} else {
						item.Payload = payload
						out <- item
					}
				}
			}()
		}
		logWaitGroup.Add(1)
		go func() {
			for logItem := range log {
				logFile.WriteString(logItem)
			}
			logWaitGroup.Done()
		}()
		dowloadWaitGroup.Wait()
		close(out)
		close(log)
		logWaitGroup.Wait()
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
