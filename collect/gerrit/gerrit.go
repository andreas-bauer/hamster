package gerrit

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/michaeldorner/hamster/collect"
)

func Generate(crawlRun collect.CrawlRun) <-chan collect.Unit {
	units := make(chan collect.Unit)
	go func() {
		defer close(units)
		startDate := crawlRun.Config.FromDate
		endDate := crawlRun.Config.ToDate

		size := int(endDate.Sub(startDate).Hours()/24) + 1
		counter := 0

		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			counter = counter + 1
			fmt.Printf("\rProcessing day %v (%v/%v - %3.1f %%)", d.Format("2006-01-02"), counter, size, float32(counter)/float32(size)*100.0)

			t1 := d.Format("2006-01-02 15:04:05.000")
			t2 := d.AddDate(0, 0, 1).Add(-1 * time.Millisecond).Format("2006-01-02 15:04:05.000")

			offset := 0
			for {
				url := fmt.Sprintf("%s/changes/?q=after:{%s}+before:{%s}&S=%v", crawlRun.Config.URL, url.QueryEscape(t1), url.QueryEscape(t2), offset)
				response_body, err := crawlRun.HTTPClient.Get(url)
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
					url := fmt.Sprintf("%s/changes/%s/detail/?o=ALL_REVISIONS&o=ALL_COMMITS&o=ALL_FILES&o=REVIEWED&o=WEB_LINKS&o=COMMIT_FOOTERS", crawlRun.Config.URL, id)
					units <- collect.Unit{
						ID: id, 
						URL: url,
					}
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
		fmt.Println("") // nice finish :)
	}()
	return units
}

func PostProcess(in <-chan collect.Unit, crawlRun collect.CrawlRun) <-chan collect.Unit {
	units := make(chan collect.Unit)
	go func() {
		defer close(units)
		for unit := range in {
			unit.Payload = unit.Payload[5:]
			units <- unit
		}
	}()
	return units
}
