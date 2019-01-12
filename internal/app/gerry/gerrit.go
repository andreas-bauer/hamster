package gerrit

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
	"strings"

	"github.com/michaeldorner/hamster/pkg/crawl"
	"github.com/michaeldorner/hamster/pkg/http"
	"github.com/michaeldorner/hamster/pkg/store"
)


var Feed crawl.Feed = func(options crawl.Options, client http.Client, repository store.Repository) <-chan crawl.Item {
	items := make(chan crawl.Item)
	go func() {
		defer close(items)

		fmt.Println("Check available parameters")

		firstChange := getFirstChange(options.URL, client)
		baseOptionsDetail := getAvailableOptions(fmt.Sprintf("%s/changes/%v/detail/?", options.URL, firstChange["_number"]), client)
		baseOptionsQuery := getAvailableOptions(fmt.Sprintf("%s/changes/?q=change:%v&", options.URL, firstChange["_number"]), client)

		fmt.Println("Create time frames")

		timeframes := crawl.GenerateTimeFrames(options.Period)

		fmt.Println("Start crawling")

		size := len(timeframes)
		start := time.Now()

		for i, timeframe := range timeframes {
			S := 0
			for {
				from := url.QueryEscape(timeframe.From.String())
				to := url.QueryEscape(timeframe.To.String())
				url := fmt.Sprintf("%s/changes/?q=after:{%s}+before:{%s}&S=%v", options.URL, from, to, S)
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
				for _, change := range jsonResponse {
					id := fmt.Sprintf("%v", change["_number"])

					urlOptions := baseOptionsDetail

					if changeHasRevision(id, options.URL, client) {
						urlOptions += "&o=ALL_REVISIONS"
					}
					url := fmt.Sprintf("%s/changes/%s/detail/?%s", options.URL, id, urlOptions)
					items <- crawl.Item{
						ID:  id+ "_d",
						URL: url,
						FileNameExtensions: "json",
					}

					url_query := fmt.Sprintf("%s/changes/?q=change:%s&%s", options.URL, id, baseOptionsQuery)
					items <- crawl.Item{
						ID:  id + "_q",
						URL: url_query,
						FileNameExtensions: "json",
					}

					_, exists := change["_more_changes"]
					more = more || exists
				}

				if more {
					S = S + len(jsonResponse)
				} else {
					break
				}
			}
			elapsed_time := time.Since(start)
			progress := float64(i+1) / float64(size)
			remaining_time := time.Duration(elapsed_time.Seconds() / progress * float64(time.Second))

			fmt.Printf("\r\f%v/%v (%.2f %%) [%v | %v]", i+1, size, progress * 100.0, elapsed_time.Round(time.Second), remaining_time.Round(time.Second))
		}
		fmt.Println("") // nice finish :)
	}()
	return items
}

func changeHasRevision(id string, baseURL string, client http.Client) bool {
	url := fmt.Sprintf("%s/changes/?q=%s&o=CURRENT_REVISION", baseURL, id)
	response_body, _ := client.Get(url)
	return string(response_body) != ")]}'\n[]\n"
}

func getFirstChange(baseURL string, client http.Client) map[string]interface{} {
	url := fmt.Sprintf("%s/changes/", baseURL)
	response_body, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	jsonResponse := make([]map[string]interface{}, 0)
	err = json.Unmarshal(response_body[5:], &jsonResponse)
	if err != nil {
		panic(err)
	}
	return jsonResponse[0]
}

func getAvailableOptions(url string, client http.Client) string {
	availableOptions := []string{}
	for _, option := range []string{"CHECK", "DOWNLOAD_COMMANDS", "ALL_COMMITS", "ALL_FILES", "WEB_LINKS", "COMMIT_FOOTERS", "LABELS", "DETAILED_LABELS", "DETAILED_ACCOUNTS", "REVIEWER_UPDATES", "MESSAGES"} {
		httpStatus, err := client.GetHTTPStatus(url+"o="+option)
		if err != nil {
			panic(err)
		}

		if httpStatus == 200 {
			availableOptions = append(availableOptions, "o=" + option)
		}
	}
	return strings.Join(availableOptions, "&")
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
