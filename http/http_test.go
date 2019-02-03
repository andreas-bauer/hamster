package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient(60, 5, nil)
	if c.maxRetries != 5 || c.hc.Timeout != time.Duration(60)*time.Second {
		t.Error("could not create new client")
	}
}

func TestGetHTTPStatus(t *testing.T) {
	mockHTTP := NewTestClient(func(req *http.Request) *http.Response {
		reqStatusCode := strings.Split(req.URL.String(), "status/")
		resStatusCode, _ := strconv.Atoi(reqStatusCode[1])

		return &http.Response{
			StatusCode: resStatusCode,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Some Payload")),
			Header:     make(http.Header),
		}
	})

	c := Client{hc: *mockHTTP, maxRetries: 1, logChan: nil}
	statusCodes := []int{200, 201, 202, 203, 204, 205, 206, 304, 307, 308, 400, 401, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 421, 426, 428, 429, 430, 431, 451, 500, 501, 502, 503, 504, 505} // 100, 101, 301, 302, 303 not tested
	var wg sync.WaitGroup
	for _, statusCode := range statusCodes {
		wg.Add(1)
		func(actualStatusCode int) {
			url := fmt.Sprintf("https://mock/status/%v", actualStatusCode)
			statusCode := c.GetHTTPStatus(url)
			if statusCode != actualStatusCode {
				t.Errorf("Expected status code %v, got %v\n", actualStatusCode, statusCode)
			}
			wg.Done()
		}(statusCode)
	}
	wg.Wait()
}

func TestPayload(t *testing.T) {
	mockHTTP := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Some Payload")),
			Header:     make(http.Header),
		}
	})

	c := Client{hc: *mockHTTP, maxRetries: 1, logChan: nil}
	response := c.Get("https://mock/anything/123")

	if len(response.Payload) != 12 {
		t.Errorf("Expected length of payload %v, got %v\n", 12, len(response.Payload))
	}
}

func TestRetry(t *testing.T) {
	mockHTTP := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 501,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`Server Error`)),
			Header:     http.Header{"Retry-After": []string{"1"}},
		}
	})

	retries := uint(1)
	c := Client{hc: *mockHTTP, maxRetries: 1, logChan: nil}
	response := c.Get("https://mock/status/501")

	if response.Retries != retries {
		t.Errorf("Expected %v retries, got %v\n", retries, response.Retries)
	}
}

func TestLog(t *testing.T) {
	mockHTTP := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
			Header:     make(http.Header),
		}
	})

	url := "https://mock/anything/123"
	log := make(chan LogEntry)
	c := Client{hc: *mockHTTP, maxRetries: 1, logChan: log}
	go c.Get(url)
	l := <-log
	close(log)

	if url != l.URL {
		t.Errorf("Expected %v, got %v\n", url, l.URL)
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
