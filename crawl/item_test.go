package crawl

import (
	"github.com/michaeldorner/hamster/http"
	"testing"
)

func TestFileName(t *testing.T) {
	item := Item{
		ID:                 "0001",
		URL:                "https://www.google.com",
		Response:           http.Response{},
		FileNameExtensions: "json",
	}
	if item.FileName() != "0001.json" {
		t.Errorf("Expected %v, got %v\n", "0001.json", item.FileName())
	}
}
