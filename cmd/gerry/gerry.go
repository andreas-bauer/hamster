package main

import (
	"flag"
	"fmt"
	"github.com/michaeldorner/hamster/internal/app/gerry"
	"github.com/michaeldorner/hamster/pkg/crawl"
)

func main() {
	configFile := flag.Arg(0)
	flag.Parse()

	if flag.NFlag() != 1 {
		panic("no configuration file ")
	}
	fmt.Println("Loading ", configFile)

	options := crawl.LoadOptionsFromJSONFile(configFile)

	if options.MaxRetryAttempts == 0 {
		options.MaxRetryAttempts = 5
	}
	if options.Timeout == 0 {
		options.Timeout == 120
	}

	crawl.Run(options, gerrit.Feed, gerrit.PostProcess)
}
