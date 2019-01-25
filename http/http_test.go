package http

import (
	"fmt"
	"strings"
	"testing"
	"sync"
)


func TestGetHTTPStatus(t *testing.T) {
	log := make(chan string)
	defer close(log)
	c := NewClient(60, 1, log)

	statusCodes := []int{200, 201, 202, 203, 204, 205, 206, 301, 302, 303, 304, 307, 308, 400, 401, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 421, 426, 428, 429, 430, 431, 451, 500, 501, 502, 503, 504, 505} // 100, 101 not tested
	var wg sync.WaitGroup
	for _, statusCode := range statusCodes {
		wg.Add(1)
		go func(actualStatusCode int) {
			url := fmt.Sprintf("https://httpbin.org/status/%v", actualStatusCode)
			statusCode, err := c.GetHTTPStatus(url)
			if !(statusCode == actualStatusCode || statusCode == 200) || err != nil {
				t.Errorf("Expected status code %v, got %v and error %v\n", actualStatusCode, statusCode, err)
			}
			wg.Done()
		}(statusCode)
	}
	wg.Wait()
}

func TestGet(t *testing.T) {
	log := make(chan string, 1) // expecting exactly one logItem

	c := NewClient(5, 1, log)
	data, err := c.Get("https://httpbin.org/anything/123")

	if err != nil && len(data) > 0 {
		t.Errorf("Client error %v\n", err)
	}	

	close(log)

	logItem := <- log // expecting exactly one logItem
	if !strings.Contains(logItem, "https://httpbin.org/anything/123") {
		t.Errorf("Expected to contain 'https://httpbin.org/anything/123', but is %v\n", logItem)
	} 
}
