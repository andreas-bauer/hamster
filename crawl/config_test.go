package crawl

import (
	"bytes"
	"encoding/json"
	"testing"
)

var configurationJSONData = `{
	"url": "https://android-review.googlesource.com",
	"period": {
		"from": "2008-07-01 00:00:00.000", 
		"to": "2018-12-31 00:00:00.000",
		"chunkSize": "24h0m0s"
	},
	"outDir": "./android/",
	"maxRetryAttempts": 10,
	"timeout": 120,
	"skipExistingFiles": false,
	"parallelRequests": 2
}`

func TestUnmarshal(t *testing.T) {
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
		t.Error("Expecting '120' for Timeout, got", configuration.Timeout)
	}

	if configuration.SkipExistingFiles {
		t.Error("Expecting 'false' for SkipExistingFiles, got 'true'")
	}

	if configuration.ParallelRequests != 2 {
		t.Error("Expecting 2 for ParallelRequests, got", configuration.ParallelRequests)
	}
}

func TestUnmarshalConfiguration(t *testing.T) {
	_, err := UnmarshalConfiguration([]byte(configurationJSONData))
	if err != nil {
		t.Error("Configuration Unmarshal error", err)
	}
}

func TestJSON(t *testing.T) {
	configuration := Configuration{}
	err := json.Unmarshal([]byte(configurationJSONData), &configuration)
	if err != nil {
		t.Error("JSON unmarshal error", err)
	}

	expectedCompactedBuffer := new(bytes.Buffer)
	err = json.Compact(expectedCompactedBuffer, []byte(configurationJSONData))
	if err != nil {
		t.Error("JSON compact error", err)
	}

	compactedBuffer := new(bytes.Buffer)
	err = json.Compact(compactedBuffer, configuration.JSON())
	if err != nil {
		t.Error("JSON compact error", err)
	}
	if !bytes.Equal(expectedCompactedBuffer.Bytes(), compactedBuffer.Bytes()) {
		t.Errorf("Expected %v, got %v\n", expectedCompactedBuffer, compactedBuffer)
	}
}
