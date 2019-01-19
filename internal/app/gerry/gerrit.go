package gerrit

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/michaeldorner/hamster/pkg/crawl"
	"github.com/michaeldorner/hamster/pkg/http"
	"github.com/michaeldorner/hamster/pkg/store"
)

type ChangeInfo struct {
	Number      int  `json:"_number"`
	MoreChanges bool `json:"_more_changes"`
}

var Feed crawl.Feed = func(configuration crawl.Configuration, client http.Client, repository store.Repository) <-chan crawl.Item {
	items := make(chan crawl.Item)
	go func() {
		defer close(items)

		fmt.Println("Check available parameters")

		firstChange := getFirstChange(configuration.URL, client)
		baseConfigurationDetail := getAvailableConfiguration(fmt.Sprintf("%s/changes/%v/detail/?", configuration.URL, firstChange.Number), client)
		baseConfigurationQuery := getAvailableConfiguration(fmt.Sprintf("%s/changes/?q=change:%v&", configuration.URL, firstChange.Number), client)

		fmt.Println("Create time frames")

		timeframes := crawl.GenerateTimeFrames(configuration.Period)

		fmt.Println("Start crawling")

		size := len(timeframes)
		start := time.Now()

		for i, timeframe := range timeframes {
			S := 0
			for {
				from := url.QueryEscape(timeframe.From.String())
				to := url.QueryEscape(timeframe.To.String())
				url := fmt.Sprintf("%s/changes/?q=after:{%s}+before:{%s}&S=%v", configuration.URL, from, to, S)
				response_body, err := client.Get(url, nil)
				if err != nil {
					panic(err)
				}
				changes := make([]ChangeInfo, 0)
				err = json.Unmarshal(response_body[5:], &changes)
				if err != nil {
					panic(err)
				}

				more := false
				for _, change := range changes {
					id := fmt.Sprintf("%v", change.Number)

					urlConfiguration := baseConfigurationDetail

					if changeHasRevision(id, configuration.URL, client) {
						urlConfiguration += "&o=ALL_REVISIONS"
					}
					url := fmt.Sprintf("%s/changes/%s/detail/?%s", configuration.URL, id, urlConfiguration)
					items <- crawl.Item{
						ID:                 id + "_d",
						URL:                url,
						FileNameExtensions: "json",
					}

					url_query := fmt.Sprintf("%s/changes/?q=change:%s&%s", configuration.URL, id, baseConfigurationQuery)
					items <- crawl.Item{
						ID:                 id + "_q",
						URL:                url_query,
						FileNameExtensions: "json",
					}

					more = more || change.MoreChanges
				}

				if more {
					S = S + len(changes)
				} else {
					break
				}
			}
			elapsed_time := time.Since(start)
			progress := float64(i+1) / float64(size)
			remaining_time := time.Duration(elapsed_time.Seconds()/progress*float64(time.Second)) - elapsed_time

			fmt.Printf("\r%v/%v (%.2f %%) [%v | %v]", i+1, size, progress*100.0, elapsed_time.Round(time.Second), remaining_time.Round(time.Second))
		}
		fmt.Println("") // nice finish :)
	}()
	return items
}

func changeHasRevision(id string, baseURL string, client http.Client) bool {
	url := fmt.Sprintf("%s/changes/?q=%s&o=CURRENT_REVISION", baseURL, id)
	response_body, _ := client.Get(url, nil)
	return string(response_body) != ")]}'\n[]\n"
}

func getFirstChange(baseURL string, client http.Client) ChangeInfo {
	url := fmt.Sprintf("%s/changes/?n=1", baseURL)
	response_body, err := client.Get(url, nil)
	if err != nil {
		panic(err)
	}
	jsonResponse := make([]ChangeInfo, 0)
	err = json.Unmarshal(response_body[5:], &jsonResponse)
	if err != nil {
		panic(err)
	}
	return jsonResponse[0]
}

func getAvailableConfiguration(url string, client http.Client) string {
	availableConfiguration := []string{}
	for _, option := range []string{"CHECK", "DOWNLOAD_COMMANDS", "ALL_COMMITS", "ALL_FILES", "WEB_LINKS", "COMMIT_FOOTERS", "LABELS", "DETAILED_LABELS", "DETAILED_ACCOUNTS", "REVIEWER_UPDATES", "MESSAGES"} {
		httpStatus, err := client.GetHTTPStatus(url + "o=" + option)
		if err != nil {
			panic(err)
		}

		if httpStatus == 200 {
			availableConfiguration = append(availableConfiguration, "o="+option)
		}
	}
	return strings.Join(availableConfiguration, "&")
}

var PostProcess crawl.PostProcess = func(configuration crawl.Configuration, client http.Client, in <-chan crawl.Item) <-chan crawl.Item {
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
