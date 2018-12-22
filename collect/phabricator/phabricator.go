package phabricator

import (
	"fmt"
	"github.com/michaeldorner/hamster/collect"
)

func Generate(crawlRun collect.CrawlRun) <-chan Unit {
	units := make(chan Unit)
	defer close(units)
	fmt.Println("Phabricator is not supported yet")
	return units
}

func PostProcessPhabricatorUnits(<-chan Unit) <-chan Unit {
	units := make(chan Unit)
	defer close(units)
	return units
}
