package crawl

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalConfiguration(t *testing.T) {
	var configurationJSONData = `{
		"url":"https://android-review.googlesource.com",
		"period": {
			"from": "2008-07-01 00:00:00.000", 
			"to":"2018-12-31 00:00:00.000",
			"chunkSize": "24h"
		},
		"outDir":"./android/",
		"maxRetryAttempts":10,
		"timeout":120,
		"skipExistingFiles":false,
		"parallelRequests":2
	}`
	configuration := Configuration{}
	err := json.Unmarshal([]byte(configurationJSONData), &configuration)
	if err != nil {
		t.Error("JSON unmarshal error", err)
	}
	if configuration.URL != "https://android-review.googlesource.com" {
		t.Error("Expecting 'https://android-review.googlesource.com', got", configuration.URL)
	}

	if configuration.MaxRetryAttempts != 10 {
		t.Error("Expecting '10' for MaxRetryAttempts, got ", configuration.MaxRetryAttempts)
	}

	if configuration.Timeout != 120 {
		t.Error("Expecting '120' for Timeout, got ", configuration.Timeout)
	}

	if configuration.SkipExistingFiles {
		t.Error("Expecting 'false' for SkipExistingFiles, got 'true'")
	}
}
