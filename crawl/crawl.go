package crawl

import (
	"sync"
	"time"
	"fmt"

	"github.com/michaeldorner/hamster/http"
	"github.com/michaeldorner/hamster/store"
)

type Feed func(Configuration, http.Client, store.Repository) <-chan Item
type PostProcess func(Configuration, http.Client, <-chan Item) <-chan Item

func Run(config Configuration, feed Feed, postProcess PostProcess) {
	repository := store.NewRepository(config.OutDir)

	log := make(chan http.ResponseMeta)
	
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		logFile := repository.LogFile()
		for r := range log {
			timestamp := time.Now()
			str := fmt.Sprintf("%v\t%v\t%v\t%v\n", timestamp.Format(time.RFC3339), r.StatusCode, r.URL, r.After.String())
			logFile.WriteString(str)
		}
		wg.Done()
	}()

	client := http.NewClient(config.Timeout, config.MaxRetryAttempts, log)

	storeConfiguration(config, repository)
	afterFeed := feed(config, client, repository)
	afterFilter := filter(config, repository, afterFeed)
	afterPayload := getPayload(client, afterFilter, config.ParallelRequests)
	afterPostProcess := postProcess(config, client, afterPayload)
	afterPersist := persist(repository, afterPostProcess)

	<-afterPersist
	close(log)
	wg.Wait()
}

func storeConfiguration(config Configuration, repository store.Repository) {
	jsonData, jsonErr := config.JSON()
	if jsonErr != nil {
		panic(jsonErr)
	}
	path := repository.ConfigurationFilePath()
	storeErr := repository.Store(path, jsonData)
	if storeErr != nil {
		panic(storeErr)
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
					response := client.Get(item.URL)
					if response.StatusCode == 200 {
						item.Payload = response.Payload
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
