package crawl

import (
	"sync"

	"github.com/michaeldorner/hamster/http"
	"github.com/michaeldorner/hamster/store"
)

type Feed func(Configuration, http.Client, store.Repository) <-chan Item
type PostProcess func(Configuration, http.Client, <-chan Item) <-chan Item

func Run(config Configuration, feed Feed, postProcess PostProcess) {
	repository := store.NewRepository(config.OutDir)

	log := make(chan string) 
	logFile := repository.LogFile()
	go func() {
		for logItem := range log {
			logFile.WriteString(logItem)
		}
	}()

	client := http.NewClient(config.Timeout, config.MaxRetryAttempts, log)

	storeConfiguration(config, repository)
	afterFeed := feed(config, client, repository)
	afterFilter := filter(config, repository, afterFeed)
	afterPayload := getPayload(client, afterFilter, config.ParallelRequests)
	afterPostProcess := postProcess(config, client, afterPayload)
	afterPersist := persist(repository, afterPostProcess)

	<- afterPersist
	close(log)
}

func storeConfiguration(config Configuration, repository store.Repository) {
	jsonData := config.JSON()
	path := repository.ConfigurationFilePath()
	err := repository.Store(path, jsonData)
	if err != nil {
		panic(err)
	}
}

func filter(config Configuration, repository store.Repository, in <-chan Item) <-chan Item {
	out := make(chan Item)
	go func() {
		defer close(out)
		for item := range in {
			path := repository.AppendDataPath(item.FileName())
			if !(config.SkipExistingFiles && repository.FileExists(path)) {
				out <- item
			}
		}
	}()
	return out
}

func getPayload(client http.Client, in <-chan Item, numParallelRequests uint) <-chan Item {
	out := make(chan Item)

	go func() {
		var parallelWaitGroup sync.WaitGroup

		for i := uint(0); i < numParallelRequests; i++ {
			parallelWaitGroup.Add(1)
			go func() {
				defer parallelWaitGroup.Done()
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
		}
		parallelWaitGroup.Wait()
		close(out)
	}()
	return out
}

func persist(repository store.Repository, in <-chan Item) chan bool {
	done := make(chan bool)
	go func() {
		for item := range in {
			file_path := repository.AppendDataPath(item.FileName())
			err := repository.Store(file_path, item.Payload)
			if err != nil {
				panic(err)
			}
		}
		close(done)
	}()
	return done
}
