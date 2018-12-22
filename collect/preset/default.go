package preset

import (
	"github.com/michaeldorner/hamster/collect"
)

func Filter(in <-chan Unit, crawlRun CrawlRun) <-chan Unit {
	out := make(chan Unit)
	go func() {
		defer close(out)
		for unit := range in {
			if crawlRun.configuration.SkipExistingFiles != crawlRun.persistence.UnitFileExists(unit.ID) {
				out <- unit
			}
		}
	}()
	return out
}

func GetPayload(in <-chan Unit, crawlRun CrawlRun) <-chan Unit {
	out := make(chan Unit)
	go func() {
		defer close(out)
		for unit := range in {
			payload, err := crawlRun.httpClient.Get(unit.URL)
			if err != nil {
				panic(err)
			} else {
				unit.Payload = payload
				out <- unit
			}
		}
	}()
	return out
}


func Store(in <-chan Unit, crawlRun CrawlRun) {
	for unit := range in {
		err := crawlRun.persistence.StoreUnit(unit.ID, unit.Payload)
		if err != nil {
			panic(err)
		}
	}
}