package main

import (
	//"flag"

	"github.com/michaeldorner/hamster/internal/app/gerry"
	"github.com/michaeldorner/hamster/pkg/crawl"
)


func main() {
	options := crawl.TestOptions
	crawl.Run(options, gerrit.Feed, gerrit.PostProcess) 
}
