package http

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var UnexpectedPanicErr = errors.New("unexpected HTTP client panic occured")

type Client struct {
	hc         http.Client
	maxRetries uint
}

func NewClient(timeOut time.Duration, maxRetries uint) Client {
	return Client{
		hc: http.Client{
			Timeout: time.Duration(timeOut),
		},
		maxRetries: maxRetries,
	}
}

type Response struct {
	StatusCode  int
	TimeToCrawl time.Duration
	Attempts    uint
	Payload     []byte
}

func (client Client) Get(url string) Response {
	response := Response{}
	startTime := time.Now()

	retryAfter := 0
	for retry := uint(0); retry <= client.maxRetries; retry++ {
		response.Attempts = retry
		if retry > 0 {
			rand.Seed(time.Now().UnixNano())
			jitter := (rand.Intn(10) + 1) * 100
			time.Sleep(time.Duration(retryAfter)*time.Second + time.Duration(jitter)*time.Millisecond)
		}
		retryAfter = 2 << retry
		r, err := client.hc.Get(url)
		response.TimeToCrawl = time.Since(startTime)

		if err == nil {
			defer r.Body.Close()
			response.StatusCode = r.StatusCode
			if r.StatusCode == 200 {
				data, err := ioutil.ReadAll(r.Body)
				if err == nil {
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
