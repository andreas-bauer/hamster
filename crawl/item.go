package crawl

import "github.com/michaeldorner/hamster/http"

type Item struct {
	ID                 string
	URL                string
	Response           http.Response
	FileNameExtensions string
}

func (item Item) FileName() string {
	return item.ID + "." + item.FileNameExtensions
}
