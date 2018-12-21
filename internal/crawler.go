package internal

type Unit struct {
	ID      string
	URL     string
	Payload []byte
}

func NewUnit(id, url string) Unit {
	return Unit{
		ID:  id,
		URL: url,
	}
}

type CrawlRun struct {
	configuration    Configuration
	httpClient       RetryHTTPClient
	persistence      Persistence
	progress         string
	createUnits      func() <-chan Unit
	postProcessUnits func(<-chan Unit) <-chan Unit
}

func NewCrawlRun(configuration Configuration) CrawlRun {
	var cr CrawlRun = CrawlRun{
		configuration: configuration,
		httpClient:    NewRetryHTTPClient(configuration.TimeOut, configuration.Retries),
		persistence:   NewPersistence(configuration.OutDir, configuration.CrawlRunID),
	}
	switch configuration.ReviewTool {
	case Gerrit:
		cr.createUnits = cr.createGerritUnits
		cr.postProcessUnits = cr.postProcessGerritUnits
	default:
		panic(ToolNotSupported)
	}
	return cr
}

func (crawlRun CrawlRun) Start() {
	crawlRun.persistence.StoreConfiguration(crawlRun.configuration)

	channel_1 := crawlRun.createUnits()
	channel_2 := crawlRun.filter(channel_1)
	channel_3 := crawlRun.getPayload(channel_2)
	channel_4 := crawlRun.postProcessUnits(channel_3)
	crawlRun.storeToFile(channel_4)
}

func (crawlRun CrawlRun) filter(in <-chan Unit) <-chan Unit {
	out := make(chan Unit)
	go func() {
		defer close(out)
		for unit := range in {
			if crawlRun.configuration.SkipExistingFiles != crawlRun.persistence.UnitFileExists(unit.ID) {
				out <- unit
			}
		}
	}()
	return out
}

func (crawlRun CrawlRun) getPayload(in <-chan Unit) <-chan Unit {
	out := make(chan Unit)
	go func() {
		defer close(out)
		for unit := range in {
			payload, err := crawlRun.httpClient.Get(unit.URL)
			if err != nil {
				panic(err)
			} else {
				unit.Payload = payload
				out <- unit
			}
		}
	}()
	return out
}

func (crawlRun CrawlRun) storeToFile(in <-chan Unit) {
	for unit := range in {
		err := crawlRun.persistence.StoreUnit(unit.ID, unit.Payload)
		if err != nil {
			panic(err)
		}
	}
}
