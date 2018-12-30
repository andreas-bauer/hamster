package client

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ErrMaxRetries = errors.New("error reached max retries")

type HamsterClient struct {
	hc         http.Client
	maxRetries int
	log        *log.Logger
}

func New(timeOut, maxRetries int, logFile *os.File) HamsterClient {
	return HamsterClient{
		hc: http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		},
		maxRetries: maxRetries,
		log:        log.New(logFile, "", 0),
	}
}

func (client HamsterClient) Get(url string) ([]byte, error) {
	retry := 0
	for {
		wait := retry * retry
		response, _ := client.hc.Get(url)

		if response != nil {
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK && response.Body != nil {
				client.log.Println("DOWNLOADED", url)
				return ioutil.ReadAll(response.Body)
			} else {
				client.log.Println("RETRY", retry, url)

				header := response.Header.Get("Retry-After")
				if len(header) > 0 {
					parsedInt, parseErr := strconv.Atoi(header)
					if parseErr != nil {
						wait = parsedInt
					}
				}
			}
		}

		if retry <= client.maxRetries {
			time.Sleep(time.Duration(wait) * time.Second)
			retry = retry + 1
		} else {
			client.log.Println("FAILED", retry, url)
			return nil, ErrMaxRetries
		}
	}
}
