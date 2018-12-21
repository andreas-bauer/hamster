package internal

func (crawlRun CrawlRun) createPhabricatorUnits() <-chan Unit {
	units := make(chan Unit)
	defer close(units)
	return units
}

func (crawlRun CrawlRun) postProcessPhabricatorUnits(<-chan Unit) <-chan Unit {
	units := make(chan Unit)
	defer close(units)
	return units
}
