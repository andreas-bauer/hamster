package preset

import (
	"github.com/michaeldorner/hamster/collect"
)

func Filter(in <-chan collect.Unit, crawlRun collect.CrawlRun) <-chan collect.Unit {
	out := make(chan collect.Unit)
	go func() {
		defer close(out)
		for unit := range in {
			if crawlRun.SkipExistingFiles != crawlRun.Persistence.UnitFileExists(unit.ID) {
				out <- unit
			}
		}
	}()
	return out
}

func GetPayload(in <-chan collect.Unit, crawlRun collect.CrawlRun) <-chan collect.Unit {
	out := make(chan collect.Unit)
	go func() {
		defer close(out)
		for unit := range in {
			payload, err := crawlRun.HTTPClient.Get(unit.URL)
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

func Store(in <-chan collect.Unit, crawlRun collect.CrawlRun) {
	for unit := range in {
		err := crawlRun.Persistence.StoreUnit(unit.ID, unit.Payload)
		if err != nil {
			panic(err)
		}
	}
}
