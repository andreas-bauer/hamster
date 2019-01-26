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
	logChan chan ResponseMeta
}

func NewClient(timeOut, maxRetries uint, lc chan ResponseMeta) Client {
	return Client{
		hc: http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		},
		maxRetries: maxRetries,
		logChan: lc,
	}
}

type ResponseMeta struct {
	StatusCode int
	After      time.Duration
	Retries    uint
	URL        string
}

type Response struct {
	ResponseMeta
	Payload []byte
}

func (client Client) Get(url string) Response {
	response := Response{}
	response.URL = url
	response.StatusCode = 444

	defer func() { 
		if client.logChan != nil {
			client.logChan <- response.ResponseMeta
		}
	}()

	retryAfter := 0
	for retry := uint(0); retry <= client.maxRetries; retry++ {
		time.Sleep(time.Duration(retryAfter) * time.Second)
		retryAfter = 2 << retry

		startTime := time.Now()
		r, err := client.hc.Get(url)
		response.After = time.Since(startTime)
		response.Retries = retry
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
					return response
				}
			}
		}
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


/*
timestamp := time.Now()
str := fmt.Sprintf("%v\t%v\t%v\t%v\n", timestamp.Format(time.RFC3339), responseMeta.StatusCode, responseMeta.URL, responseMeta.After.String())

client.logFileMutex.Lock()
client.logFileWriter.Write([]byte(str))
client.logFileMutex.Unlock()
*/