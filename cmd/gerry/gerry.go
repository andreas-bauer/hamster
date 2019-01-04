package main

import (
	"fmt"
	"github.com/michaeldorner/hamster/internal/app/gerry"
	"github.com/michaeldorner/hamster/pkg/crawl"
	"os"
)

func main() {
	configFile := os.Args[1]

	fmt.Println("Loading ", configFile)

	options := crawl.LoadOptionsFromJSONFile(configFile)

	crawl.Run(options, gerrit.Feed, gerrit.PostProcess)
}
