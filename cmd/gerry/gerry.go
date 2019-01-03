package main

import (
	"os"
	"fmt"
	"github.com/michaeldorner/hamster/internal/app/gerry"
	"github.com/michaeldorner/hamster/pkg/crawl"
)

func main() {
	configFile := os.Args[1]

	fmt.Println("Loading ", configFile)

	options := crawl.LoadOptionsFromJSONFile(configFile)

	crawl.Run(options, gerrit.Feed, gerrit.PostProcess)
}
