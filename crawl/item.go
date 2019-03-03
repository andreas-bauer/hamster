package crawl

import (
	"github.com/michaeldorner/hamster/http"
)

type Item struct {
	ID                 string
	Request            *http.Request
	Response           *http.Response
	FileNameExtensions string
}

func NewItem(id, url, FileNameExtensions string) (*Item, error) {
	request, err := http.NewGetRequest(url)
	if err != nil {
		return nil, err
	}
	return &Item{
		ID:                 id,
		Request:            request,
		FileNameExtensions: FileNameExtensions,
	}, nil
}

func (item Item) FileName() string {
	return item.ID + "." + item.FileNameExtensions
}
