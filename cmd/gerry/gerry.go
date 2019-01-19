package main

import (
	"fmt"
	"github.com/michaeldorner/hamster/internal/app/gerry"
	"github.com/michaeldorner/hamster/pkg/crawl"
	"os"
)

func main() {
	configFile := os.Args[1]

	fmt.Println("Load", configFile)
	configuration := crawl.LoadConfigurationFromJSONFile(configFile)

	fmt.Printf("Start crawl run with %v parallel requests", configuration.ParallelRequests)
	crawl.Run(configuration, gerrit.Feed, gerrit.PostProcess)
}
