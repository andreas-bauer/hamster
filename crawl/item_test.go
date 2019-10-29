package crawl

import (
	"testing"

	"github.com/michaeldorner/hamster/http"
)

func TestFileName(t *testing.T) {
	item := Item{
		ID:                 "0001",
		Request:            &http.Request{},
		Response:           &http.Response{},
		FileNameExtensions: "json",
	}
	if item.FileName() != "0001.json" {
		t.Errorf("Expected %v, got %v\n", "0001.json", item.FileName())
	}
}
