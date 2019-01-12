package crawl

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalOptions(t *testing.T) {
	var optionsJSONData = `{
		"url":"https://android-review.googlesource.com",
		"period": {
			"from": "2008-07-01 00:00:00.000", 
			"to":"2018-12-31 00:00:00.000",
			"chunkSize": "24h"
		},
		"outDir":"./android/",
		"maxRetryAttempts":10,
		"timeout":120,
		"skipExistingFiles":false
	}`
	options := Options{}
	err := json.Unmarshal([]byte(optionsJSONData), &options)
	if err != nil {
		t.Error("JSON unmarshal error", err)
	}
	if options.URL != "https://android-review.googlesource.com" {
		t.Error("Expecting 'https://android-review.googlesource.com', got", options.URL)
	}

	if options.MaxRetryAttempts != 10 {
		t.Error("Expecting '10' for MaxRetryAttempts, got ", options.MaxRetryAttempts)
	}

	if options.Timeout != 120 {
		t.Error("Expecting '120' for Timeout, got ", options.Timeout)
	}

	if options.SkipExistingFiles {
		t.Error("Expecting 'false' for SkipExistingFiles, got 'true'")
	}
}
