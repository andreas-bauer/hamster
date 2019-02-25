package crawl

import (
	"testing"
	"time"
)

func parseDuration(value string, t *testing.T) Duration {
	duration, err := ParseDuration(value)
	if err != nil {
		t.Error("ParseTimestamp() failed with input", value)
	}
	return duration
}

func TestParseDuration(t *testing.T) {
	duration := parseDuration("24h", t)
	reference := 24 * time.Hour

	if duration.Duration != reference {
		t.Errorf("Expected %v, got %v\n", reference, duration.Duration)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	duration1d := parseDuration("24h", t)
	d := Duration{}
	d.UnmarshalJSON([]byte("\"24h\""))
	if d.Duration != duration1d.Duration {
		t.Errorf("Expected 24h, got %v\n", d.Duration)
	}
}
