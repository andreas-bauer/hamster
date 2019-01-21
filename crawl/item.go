package crawl

type Item struct {
	ID                 string
	URL                string
	Payload            []byte
	FileNameExtensions string
}

func (item Item) FileName() string {
	return item.ID + "." + item.FileNameExtensions
}
