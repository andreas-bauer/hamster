package collect

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

type RetryHTTPClient struct {
	hc         http.Client
	MaxRetries int
	Log        *log.Logger
}

func NewRetryHTTPClient(timeOut, maxRetries int, logFile *os.File) RetryHTTPClient {
	return RetryHTTPClient{
		hc: http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		},
		MaxRetries: maxRetries,
		Log:        log.New(logFile, "", 0),
	}
}

func (client RetryHTTPClient) Get(url string) ([]byte, error) {
	retry := 0
	for {
		wait := retry * retry
		response, _ := client.hc.Get(url)

		if response != nil {
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK && response.Body != nil {
				client.Log.Println("DOWNLOADED", url)
				return ioutil.ReadAll(response.Body)
			} else {
				client.Log.Println("RETRY", retry, url)

				header := response.Header.Get("Retry-After")
				if len(header) > 0 {
					parsedInt, parseErr := strconv.Atoi(header)
					if parseErr != nil {
						wait = parsedInt
					}
				}
			}
		}

		if retry <= client.MaxRetries {
			time.Sleep(time.Duration(wait) * time.Second)
			retry = retry + 1
		} else {
			client.Log.Println("FAILED", retry, url)
			return nil, ErrMaxRetries
		}
	}
}
