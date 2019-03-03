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
	c := NewClient(time.Duration(60)*time.Second, 5)
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

	c := Client{hc: *mockHTTP, maxRetries: 1}
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
	payload := "Some Payload"
	mockHTTP := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(payload)),
			Header:     make(http.Header),
		}
	})

	c := Client{hc: *mockHTTP, maxRetries: 1}
	req, err := NewGetRequest("https://mock/anything/123")
	if err != nil {
		t.Error(err)
	}

	response, err := c.Do(req)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(response.payload, []byte(payload)) != 0 {
		t.Errorf("Expected payload %v, got %v\n", payload, response.payload)
	}

	if response.TimeToCrawl() == time.Duration(0) {
		t.Errorf("Expected longer time to crawl than 0 ns\n")
	}
}

func TestGet(t *testing.T) {
	mockHTTP := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString("0")),
			Header:     make(http.Header),
		}
	})

	c := Client{hc: *mockHTTP, maxRetries: 1}
	url := "https://mock/status/200"
	resp1, err := c.Get(url)
	if err != nil {
		t.Error(err)
	}

	req, err := NewGetRequest(url)
	if err != nil {
		t.Error(err)
	}

	resp2, err := c.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp1.statusCode != resp2.statusCode {
		t.Errorf("Expected status code %v, got %v\n", resp1.statusCode, resp2.statusCode)
	}

	if resp1.retries != resp2.retries {
		t.Errorf("Expected retries %v, got %v\n", resp1.retries, resp2.retries)
	}

	if bytes.Compare(resp1.payload, resp2.payload) != 0 {
		t.Errorf("Expected status code %v, got %v\n", resp1.payload, resp2.payload)
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
	c := Client{hc: *mockHTTP, maxRetries: 1}

	req, err := NewGetRequest("https://mock/status/501")
	if err != nil {
		t.Error(err)
	}

	response, _ := c.Do(req) // 501 and error is expected

	if response.StatusCode() != 501 {
		t.Errorf("Expected status code %v , got %v\n", 501, response.StatusCode())
	}

	if response.Retries() != retries {
		t.Errorf("Expected %v retrys, got %v\n", retries, response.Retries())
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
