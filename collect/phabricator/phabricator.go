package phabricator

import (
	"fmt"
	"github.com/michaeldorner/hamster/collect"
)

func Generate(crawlRun collect.CrawlRun) <-chan collect.Unit {
	units := make(chan collect.Unit)
	defer close(units)
	fmt.Println("Phabricator is not supported yet")
	return units
}

func PostProcess(in <-chan collect.Unit, crawlRun collect.CrawlRun) <-chan collect.Unit {
	units := make(chan collect.Unit)
	defer close(units)
	return units
}
