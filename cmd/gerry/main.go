package main

import (
	"Gerry/internal"
	"fmt"
)

func main() {

	config := internal.TestConfiguration

	fmt.Println(config.CrawlRunID)
	fmt.Println(config.URL)

	crawlRun := internal.NewCrawlRun(config)
	crawlRun.Start()

}

// wipe ID
