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

	crawlTask := collect.LoadTaskFile(*pathPtr)

	channel_1 := phabricator.Generate(crawlTask)
	channel_2 := preset.Filter(channel_1, crawlTask)
	channel_3 := preset.GetPayload(channel_2, crawlTask)
	channel_4 := phabricator.PostProcess(channel_3, crawlTask)
	preset.Repository(channel_4, crawlTask)

}

// wipe ID
