package http

import (
	"errors"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var UnexpectedPanicErr = errors.New("unexpected HTTP client panic occured")
var MaxRetriesExceededErr = errors.New("max retries exceeded")

type Request = http.Request

func NewGetRequest(url string) (*Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	_req_ := Request(*req)
	return &_req_, err
}

type Client struct {
	hc         http.Client
	maxRetries uint
}

func NewClient(timeOut time.Duration, maxRetries uint) *Client {
	return &Client{
		hc: http.Client{
			Timeout: time.Duration(timeOut),
		},
		maxRetries: maxRetries,
	}
}

type Response struct {
	StatusCode  int
	TimeToCrawl time.Duration
	Retries     uint
	Payload     []byte
}

func (client Client) Get(url string) (*Response, error) {
	req, err := NewGetRequest(url)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func (client Client) Do(request *Request) (*Response, error) {
	response := &Response{}
	startTime := time.Now()

	for retry := uint(0); retry <= client.maxRetries; retry++ {
		response.Retries = retry
		if retry > 0 {
			rand.Seed(time.Now().UnixNano())
			jitter := time.Duration(rand.Intn(500)) * time.Millisecond
			backoff := time.Duration(math.Pow(2, float64(retry-1))) * time.Second
			time.Sleep(backoff + jitter)
		}

		r, err := client.hc.Do(request)
		response.TimeToCrawl = time.Since(startTime)

		if err == nil {
			defer r.Body.Close()
			response.StatusCode = r.StatusCode
			if r.StatusCode == 200 {
				data, err := ioutil.ReadAll(r.Body)
				if err == nil {
					response.Payload = data
					return response, err
				}
			}
		}
	}
	return response, MaxRetriesExceededErr
}

func (client Client) GetHTTPStatus(url string) int {
	response, err := client.hc.Get(url)
	if err != nil {
		return 444
	} else {
		return response.StatusCode
	}
}
