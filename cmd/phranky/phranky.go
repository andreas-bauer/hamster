package main

import (
	"flag"
	"github.com/michaeldorner/hamster/collect"
	"github.com/michaeldorner/hamster/collect/phabricator"
	"github.com/michaeldorner/hamster/collect/preset"
)

func main() {
	pathPtr := flag.String("path", "./", "path to config file")
	flag.Parse()

	crawlRun := collect.LoadCrawlRunFile(*pathPtr)

	channel_1 := phabricator.Generate(crawlRun)
	channel_2 := preset.Filter(channel_1, crawlRun)
	channel_3 := preset.GetPayload(channel_2, crawlRun)
	channel_4 := phabricator.PostProcess(channel_3, crawlRun)
	preset.Store(channel_4, crawlRun)

}

// wipe ID
