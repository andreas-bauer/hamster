package internal

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
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
		response, err := client.hc.Get(url)
		if err == nil && response.StatusCode == 200 {
			defer response.Body.Close()
			return ioutil.ReadAll(response.Body)
		} else {
			if retry <= client.MaxRetries {
				time.Sleep(time.Duration(retry*retry) * time.Second)
				retry = retry + 1
			} else {
				return nil, ErrMaxRetries
			}
		}
	}
}
