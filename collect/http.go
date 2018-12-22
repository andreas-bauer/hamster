package collect

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
)

var ErrMaxRetries = errors.New("error reached max retries")

type RetryHTTPClient struct {
	hc         http.Client
	MaxRetries int
}

func NewRetryHTTPClient(timeOut, maxRetries int) RetryHTTPClient {
	return RetryHTTPClient{
		hc: http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		},
		MaxRetries: maxRetries,
	}
}


func (client RetryHTTPClient) Get(url string) ([]byte, error) {
	retry := 0
	for {
		wait := retry * retry
		response, httpErr := client.hc.Get(url)

		if response != nil {
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK && response.Body != nil {
				return ioutil.ReadAll(response.Body)
			} else {
				header := response.Header.Get("Retry-After")
				if len(header)>0 {
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
			return nil, ErrMaxRetries
		}
	}
}

func retryAfterSeconds(response *http.Response, retry int) (bool, int) {
	wait := retry * retry
	if response != nil {
		if response.StatusCode == http.StatusOK && response.Body != nil {
			defer response.Body.Close()
			return ioutil.ReadAll(response.Body)
		} else {
			header := response.Header.Get("Retry-After")
			if len(header)>0 {
				parsedInt, parseErr := strconv.Atoi(header)
				if parseErr != nil {
					wait = parsedInt
				}
			} 
		}
	} 
	return true, wait
	
}
