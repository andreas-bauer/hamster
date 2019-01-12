package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ErrMaxRetries = errors.New("error reached max retries")

type Client struct {
	hc         http.Client
	maxRetries uint
	logFile    *os.File
}

func NewClient(timeOut, maxRetries uint, logFile *os.File) Client {
	return Client{
		hc: http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		},
		maxRetries: maxRetries,
		logFile:    logFile,
	}
}

func (client Client) Get(url string) ([]byte, error) {
	retryAttempt := uint(0)
	for {
		wait := 2 << uint(retryAttempt)
		response, err := client.hc.Get(url)
		if err != nil {
			client.log(timeout, 408, retryAttempt, url)
		}

		if response != nil {
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK && response.Body != nil {
				client.log(success, response.StatusCode, retryAttempt, url)
				return ioutil.ReadAll(response.Body)
			} else {
				client.log(retry, response.StatusCode, retryAttempt, url)

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
			client.log(failure, response.StatusCode, retryAttempt, url)
			return nil, ErrMaxRetries
		}
	}
}

type status string

const (
	success status = "SUCCESS"
	retry   status = "RETRY"
	failure status = "FAILURE"
	timeout status = "TIMEOUT"
)

func (client Client) log(status status, httpStatus int, retryAttempt uint, url string) {
	timestamp := time.Now()
	status_string := string(status)
	if status == retry {
		status_string = fmt.Sprintf("%v %v", status_string, retryAttempt)
	}
	str := fmt.Sprintf("%v\t%v\t%v\t%v\n", timestamp.Format(time.RFC3339), status_string, httpStatus, url)
	_, err := client.logFile.WriteString(str)
	if err != nil {
		panic(err)
	}
}
