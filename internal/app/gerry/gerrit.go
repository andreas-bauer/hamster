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

var Feed crawl.Feed = func(options crawl.Options, client http.Client, repository store.Repository) <-chan crawl.Item {
	items := make(chan crawl.Item)
	go func() {
		defer close(items)

		timeframes := crawl.GenerateTimeFrames(options.Period)

		size := len(timeframes)
		bar := progressbar.NewOptions(size, progressbar.OptionSetRenderBlankState(true), progressbar.OptionShowIts(), progressbar.OptionShowCount(), progressbar.OptionSetWidth(100))
		bar.RenderBlank()

		for _, timeframe := range timeframes {
			offset := 0
			for {
				from := url.QueryEscape(timeframe.From.String())
				to := url.QueryEscape(timeframe.To.String())
				url := fmt.Sprintf("%s/changes/?q=after:{%s}+before:{%s}&S=%v", options.URL, from, to, offset)
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
					defaultOptions := "o=CHECK&o=DOWNLOAD_COMMANDS&o=ALL_COMMITS&o=ALL_FILES&o=WEB_LINKS&o=COMMIT_FOOTERS"
					id := fmt.Sprintf("%v", response["_number"])

					if changeHasRevision(id, options.URL, client) {
						defaultOptions += "&o=ALL_REVISIONS"
					}
					url := fmt.Sprintf("%s/changes/%s/detail/?%s", options.URL, id, defaultOptions)
					items <- crawl.Item{
						ID:  "D" + id,
						URL: url,
					}

					detailsOptions := "o=LABELS&o=DETAILED_LABELS&o=DETAILED_ACCOUNTS&o=REVIEWER_UPDATES&o=MESSAGES"
					url_query := fmt.Sprintf("%s/changes/?q=%s&%s&%s", options.URL, id, defaultOptions, detailsOptions)
					items <- crawl.Item{
						ID:  "Q" + id,
						URL: url_query,
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
	return items
}

func changeHasRevision(id string, baseURL string, client http.Client) bool {
	url := fmt.Sprintf("%s/changes/?q=%v&o=CURRENT_REVISION", baseURL, id)
	response_body, _ := client.Get(url)
	return string(response_body) != ")]}'\n[]\n"
}

var PostProcess crawl.PostProcess = func(options crawl.Options, client http.Client, in <-chan crawl.Item) <-chan crawl.Item {
	items := make(chan crawl.Item)
	go func() {
		defer close(items)
		for item := range in {
			item.Payload = item.Payload[5:]
			items <- item
		}
	}()
	return items
}
