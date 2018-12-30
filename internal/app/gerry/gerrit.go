package gerrit

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/michaeldorner/hamster/pkg/crawl"
	"github.com/michaeldorner/hamster/pkg/client"
	"github.com/michaeldorner/hamster/pkg/store"
)

var Feed crawl.Feed = func (options crawl.Options, client client.HamsterClient, repository store.Repository) <-chan crawl.Unit {
	units := make(chan crawl.Unit)
	go func() {
		defer close(units)
		startDate := options.FromDate
		endDate := options.ToDate

		size := int(endDate.Sub(startDate).Hours()/24) + 1
		counter := 0

		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			counter = counter + 1
			fmt.Printf("\rProcessing day %v (%v/%v - %3.1f %%)", d.Format("2006-01-02"), counter, size, float32(counter)/float32(size)*100.0)

			t1 := d.Format("2006-01-02 15:04:05.000")
			t2 := d.AddDate(0, 0, 1).Add(-1 * time.Millisecond).Format("2006-01-02 15:04:05.000")

			offset := 0
			for {
				url := fmt.Sprintf("%s/changes/?q=after:{%s}+before:{%s}&S=%v", options.URL, url.QueryEscape(t1), url.QueryEscape(t2), offset)
				response_body, err := client.Get(url)
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
					url := fmt.Sprintf("%s/changes/%s/detail/?o=ALL_REVISIONS&o=ALL_COMMITS&o=ALL_FILES&o=REVIEWED&o=WEB_LINKS&o=COMMIT_FOOTERS", options.URL, id)
					units <- crawl.Unit{
						ID:  id,
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

var PostProcess crawl.PostProcess = func (options crawl.Options, client client.HamsterClient, in <-chan crawl.Unit) <-chan crawl.Unit {
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

