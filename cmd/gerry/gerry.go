package main

import (
	"flag"
	"github.com/michaeldorner/hamster/internal/app/gerry"
	"github.com/michaeldorner/hamster/pkg/crawl"
)

func main() {
	configFile := flag.String("configFile", "", "`path` to the JSON configuration file")
	flag.Parse()

	options := crawl.LoadOptionsFromJSONFile(*configFile)
	crawl.Run(options, gerrit.Feed, gerrit.PostProcess)
}
