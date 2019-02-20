package crawl

import (
	"fmt"
	"sync"
	"time"

	"github.com/michaeldorner/hamster/http"
	"github.com/michaeldorner/hamster/store"
)

type Feed func(Configuration, http.Client, *store.Repository) <-chan Item

func Run(config Configuration, feed Feed) error {
	repository, err := store.NewRepository(config.OutDir)
	if err != nil {
		return err
	}

	client := http.NewClient(config.Timeout.Duration, config.MaxRetries)

	jsonData, parseErr := config.JSON()
	if parseErr != nil {
		return parseErr
	}

	storeErr := repository.StoreConfiguration(jsonData)
	if storeErr != nil {
		return storeErr
	}

	afterFeed := feed(config, client, repository)
	afterPayload := get(client, afterFeed, config.ParallelRequests)
	afterPersist := persist(repository, afterPayload)
	done := log(repository, afterPersist)

	<-done
	return err
}

func storeConfiguration(config Configuration, repository *store.Repository) error {
	jsonData, jsonErr := config.JSON()
	if jsonErr != nil {
		return jsonErr
	}
	storeErr := repository.StoreConfiguration(jsonData)
	if storeErr != nil {
		return storeErr
	}
	return nil
}

func get(client http.Client, in <-chan Item, numParallelRequests uint) <-chan Item {
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
						item.Response = response
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

func persist(repository *store.Repository, in <-chan Item) <-chan Item {
	out := make(chan Item)
	go func() {
		for item := range in {
			err := repository.StoreItem(item.FileName(), item.Response.Payload)
			if err != nil {
				panic(err)
			}
			out <- item
		}
		close(out)
	}()
	return out
}

func log(repository *store.Repository, in <-chan Item) <-chan bool {
	done := make(chan bool)
	go func() {
		logFile, err := repository.LogFile()
		if err != nil {
			panic(err)
		}
		for item := range in {
			timestamp := time.Now()
			str := fmt.Sprintf("%v\t%v\t%v\t%v\n", timestamp.Format(time.RFC3339), item.Response.StatusCode, item.URL, item.Response.TimeToCrawl.String())
			logFile.WriteString(str)
		}
		close(done)
	}()
	return done
}
