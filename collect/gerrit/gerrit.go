package gerrit

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
	"github.com/michaeldorner/hamster/collect"
)

func Run(filePath string) {
	crawlRun, err := LoadCrawlRunFile(filePath)
	if err != nil {
		panic(err)
	}

	channel_1 := Generate(crawlRun)
	channel_2 := collect.Filter(channel_1, crawlRun)
	channel_3 := collect.GetPayload(channel_2, crawlRun)
	channel_4 := PostProcess(channel_3, crawlRun)
	
	collect.store(channel_4, crawlRun)

}

func Generate(crawlRun collect.CrawlRun) <-chan collect.Unit {
	units := make(chan collect.Unit)
	go func() {
		defer close(units)
		startDate := crawlRun.configuration.FromDate
		endDate := crawlRun.configuration.ToDate

		size := int(startDate.Sub(endDate).Hours()/24) + 1
		counter := 0

		for d := startDate; d.Before(endDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
			fmt.Println(counter / size)

			t1 := d.Format("2006-01-02 15:04:05.000")
			t2 := d.AddDate(0, 0, 1).Add(-1 * time.Millisecond).Format("2006-01-02 15:04:05.000")

			offset := 0
			for {
				url := fmt.Sprintf("%s/changes/?q=after:{%s}+before:{%s}&S=%v", crawlRun.configuration.URL, url.QueryEscape(t1), url.QueryEscape(t2), offset)
				response_body, err := crawlRun.httpClient.Get(url)
				if err != nil {
					panic(err)
				}
				jsonResponse := make([]map[string]interface{}, 0)
				err = json.Unmarshal(response_body[5:], &jsonResponse)
				if err != nil {
					panic(err)
				}

				for _, response := range jsonResponse {
					id := fmt.Sprintf("%v", response["_number"])
					url := fmt.Sprintf("%s/changes/%s/detail/?o=ALL_REVISIONS&o=ALL_COMMITS&o=ALL_FILES&o=REVIEWED&o=WEB_LINKS&o=COMMIT_FOOTERS", crawlRun.configuration.URL, id)
					units <- NewUnit(id, url)
				}
				l := len(jsonResponse)
				last := jsonResponse[l-1]
				if _, ok := last["_more_changes"]; ok {
					offset = offset + l
				} else {
					break
				}
			}
		}
	}()
	return units
}

func PostProcess(in <-chan Unit, crawlRun collect.CrawlRun) <-chan Unit {
	units := make(chan Unit)
	go func() {
		defer close(units)
		for unit := range in {
			unit.Payload = unit.Payload[5:]
			units <- unit
		}
	}()
	return units
}
