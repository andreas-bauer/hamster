package http

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	hc         http.Client
	maxRetries uint
	logChan    chan LogEntry
}

func NewClient(timeOut, maxRetries uint, lc chan LogEntry) Client {
	return Client{
		hc: http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		},
		maxRetries: maxRetries,
		logChan:    lc,
	}
}

type LogEntry struct {
	StatusCode int
	After      time.Duration
	Retries    uint
	URL        string
}

type Response struct {
	LogEntry
	Payload []byte
}

func (client Client) Get(url string) Response {
	/*
		log := func(r Response) {
			if client.logChan != nil {
				client.logChan <- r.ResponseMeta
			}
		}*/

	response := Response{}
	response.URL = url
	response.StatusCode = 444

	retryAfter := 0
	for retry := uint(0); retry <= client.maxRetries; retry++ {
		response.Retries = retry

		time.Sleep(time.Duration(retryAfter) * time.Second)
		retryAfter = 2 << retry
		startTime := time.Now()
		r, err := client.hc.Get(url)
		response.After = time.Since(startTime)

		if err == nil {
			response.StatusCode = r.StatusCode

			header := r.Header.Get("Retry-After")
			if len(header) > 0 {
				parsedInt, parseErr := strconv.Atoi(header)
				if parseErr != nil {
					retryAfter = parsedInt
				}
			}
			defer r.Body.Close()
			if r.StatusCode == 200 {
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				} else {
					response.Payload = data
					if client.logChan != nil {
						go func() {
							client.logChan <- response.LogEntry
						}()
					}
					return response
				}
			}
		}
	}
	if client.logChan != nil {
		go func() {
			client.logChan <- response.LogEntry
		}()
	}
	return response
}

func (client Client) GetHTTPStatus(url string) int {
	response, err := client.hc.Get(url)
	if err != nil {
		return 444
	} else {
		return response.StatusCode
	}
}
