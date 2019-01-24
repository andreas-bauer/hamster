package http

import (
	"fmt"
	"testing"
)

func TestGetHTTPStatus(t *testing.T) {
	log := make(chan string)
	defer close(log)
	c := NewClient(60, 1, log)

	statusCodes := []int{200, 201, 202, 203, 204, 205, 206, 301, 302, 303, 304, 307, 308, 400, 401, 403, 404, 405, 406, 407, 408, 409, 410, 411, 412, 413, 414, 415, 416, 417, 421, 426, 428, 429, 430, 431, 451, 500, 501, 502, 503, 504, 505} // 100, 101 not tested

	for _, statusCode := range statusCodes {
		url := fmt.Sprintf("https://httpbin.org/status/%v", statusCode)
		sc, err := c.GetHTTPStatus(url)
		if !(sc == statusCode || sc == 200) || err != nil {
			t.Errorf("Expected status code %v, got %v and error %v\n", statusCode, sc, err)
		}
	}
}
