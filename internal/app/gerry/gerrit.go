package gerrit

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/michaeldorner/hamster/pkg/crawl"
	"github.com/michaeldorner/hamster/pkg/http"
	"github.com/michaeldorner/hamster/pkg/store"
	"github.com/schollz/progressbar"
)

var Feed crawl.Feed = func(options crawl.Options, client http.Client, repository store.Repository) <-chan crawl.Unit {
	units := make(chan crawl.Unit)
	go func() {
		defer close(units)

		crawlRange := crawl.GenerateCrawlRange(options.FromDate, options.ToDate)

		size := len(crawlRange)
		bar := progressbar.NewOptions(size, progressbar.OptionSetRenderBlankState(true), progressbar.OptionShowIts(), progressbar.OptionShowCount(), progressbar.OptionSetWidth(100))
		bar.RenderBlank()

		for _, d := range crawlRange {

			offset := 0
			for {
				url := fmt.Sprintf("%s/changes/?q=after:{%s}+before:{%s}&S=%v", options.URL, url.QueryEscape(d.Min()), url.QueryEscape(d.Max()), offset)
				response_body, err := client.Get(url)
				if err != nil {
					panic(err)
				}
				jsonResponse := make([]map[string]interface{}, 0)
				err = json.Unmarshal(response_body[5:], &jsonResponse)
				if err != nil {
					panic(err)
				}

				more := false

				for _, response := range jsonResponse {
					id := fmt.Sprintf("%v", response["_number"])
					url := fmt.Sprintf("%s/changes/%s/detail/?o=ALL_REVISIONS&o=ALL_COMMITS&o=ALL_FILES&o=REVIEWED&o=WEB_LINKS&o=COMMIT_FOOTERS", options.URL, id)
					units <- crawl.Unit{
						ID:  id,
						URL: url,
					}
					_, exists := response["_more_changes"]
					more = more || exists
				}

				if more {
					offset = offset + len(jsonResponse)
				} else {
					break
				}
			}
			bar.Add(1)
		}
		fmt.Println("") // nice finish :)
	}()
	return units
}

var PostProcess crawl.PostProcess = func(options crawl.Options, client http.Client, in <-chan crawl.Unit) <-chan crawl.Unit {
	units := make(chan crawl.Unit)
	go func() {
		defer close(units)
		for unit := range in {
			unit.Payload = unit.Payload[5:]
			units <- unit
		}
	}()
	return units
}
