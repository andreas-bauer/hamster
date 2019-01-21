package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var ErrMaxRetries = errors.New("error reached max retries")

type Client struct {
	hc         http.Client
	maxRetries uint
}

func NewClient(timeOut, maxRetries uint) Client {
	return Client{
		hc: http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		},
		maxRetries: maxRetries,
	}
}

func (client Client) Get(url string, logChan chan string) ([]byte, error) {
	retryAttempt := uint(0)
	startTime := time.Now()
	for {
		wait := 2 << uint(retryAttempt)
		response, err := client.hc.Get(url)
		if err != nil {
			logToChan(logChan, timeout, 408, retryAttempt, url, startTime)
		} else {
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK && response.Body != nil {
				logToChan(logChan, success, response.StatusCode, retryAttempt, url, startTime)
				return ioutil.ReadAll(response.Body)
			} else {
				logToChan(logChan, retry, response.StatusCode, retryAttempt, url, startTime)

				header := response.Header.Get("Retry-After")
				if len(header) > 0 {
					parsedInt, parseErr := strconv.Atoi(header)
					if parseErr != nil {
						wait = parsedInt
					}
				}
			}
		}

		if retryAttempt <= client.maxRetries {
			time.Sleep(time.Duration(wait) * time.Second)
			retryAttempt = retryAttempt + 1
		} else {
			logToChan(logChan, failure, response.StatusCode, retryAttempt, url, startTime)
			return []byte{}, ErrMaxRetries
		}
	}
}


func (client Client) GetHTTPStatus(url string) (int, error) {
	response, err := client.hc.Get(url)
	if err != nil {
		return 0, err
	} else {
		return response.StatusCode, err
	}
}

type status string

const (
	success status = "SUCCESS"
	retry   status = "RETRY"
	failure status = "FAILURE"
	timeout status = "TIMEOUT"
)

func logToChan(logChan chan string, status status, httpStatus int, retryAttempt uint, url string, start time.Time) {
	if logChan != nil {
		timestamp := time.Now()
		status_string := string(status)
		if status == retry {
			status_string = fmt.Sprintf("%v %v", status_string, retryAttempt)
		}
		str := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", timestamp.Format(time.RFC3339), status_string, httpStatus, url, time.Since(start).String())
		logChan <- str
	}	
}
