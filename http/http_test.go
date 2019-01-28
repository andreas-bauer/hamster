package http

import (
	"fmt"
	"sync"
	"testing"
)

func TestGetHTTPStatus(t *testing.T) {
	c := NewClient(60, 1, nil)
	statusCodes := []int{200, 201, 202, 203, 204, 205, 206, 301, 302, 303, 304, 307, 308, 400, 401, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 421, 426, 428, 429, 430, 431, 451, 500, 501, 502, 503, 504, 505} // 100, 101 not tested
	var wg sync.WaitGroup
	for _, statusCode := range statusCodes {
		wg.Add(1)
		func(actualStatusCode int) {
			url := fmt.Sprintf("https://httpbin.org/status/%v", actualStatusCode)
			statusCode := c.GetHTTPStatus(url)
			if !(statusCode == actualStatusCode || statusCode == 200) {
				t.Errorf("Expected status code %v, got %v\n", actualStatusCode, statusCode)
			}
			wg.Done()
		}(statusCode)
	}
	wg.Wait()
}

func TestPayload(t *testing.T) {
	c := NewClient(5, 1, nil)
	response := c.Get("https://httpbin.org/anything/123")
	if len(response.Payload) < 1 {
		t.Errorf("no expected payload for 'https://httpbin.org/anything/123'")
	}
}

func TestRetry(t *testing.T) {
	retries := uint(1)
	c := NewClient(1, retries, nil)
	response := c.Get("https://httpbin.org/status/501")

	if response.Retries != retries {
		t.Errorf("Expected %v retries, got %v\n", retries, response.Retries)
	}
}

func TestLog(t *testing.T) {
	url := "https://httpbin.org/anything/123"
	log := make(chan LogEntry)
	c := NewClient(5, 1, log)
	c.Get(url)
	l := <-log
	close(log)

	if url != l.URL {
		t.Errorf("Expected %v, got %v\n", url, l.URL)
	}
}
