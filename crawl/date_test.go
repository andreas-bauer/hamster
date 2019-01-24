package crawl

import (
	"encoding/json"
	"testing"
	"time"
)

func parseTimestamp(value string, t *testing.T) Timestamp {
	ts, err := ParseTimestamp(value)
	if err != nil {
		t.Error("ParseTimestamp() failed with input", value)
	}
	return ts
}

func parseDuration(value string, t *testing.T) Duration {
	duration, err := ParseDuration(value)
	if err != nil {
		t.Error("ParseTimestamp() failed with input", value)
	}
	return duration
}

func TestParseTimestamp(t *testing.T) {
	timestamp := parseTimestamp("2018-01-01 00:00:00.000", t)
	reference := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)

	if !timestamp.Time.Equal(reference) {
		t.Errorf("Expected %v, got %v\n", reference, timestamp.Time)
	}
}

func TestParseDuration(t *testing.T) {
	duration := parseDuration("24h", t)
	reference := 24 * time.Hour

	if duration.Duration != reference {
		t.Errorf("Expected %v, got %v\n", reference, duration.Duration)
	}
}

func TestPeriod(t *testing.T) {
	var rangeJSONData = `
	{
		"from":"2018-01-01 00:00:00.000", 
		"to":"2018-02-01 00:00:00.000",
		"chunkSize":"24h"
	}`
	rd := Period{}
	err := json.Unmarshal([]byte(rangeJSONData), &rd)

	fromTimestamp, err := ParseTimestamp("2018-01-01 00:00:00.000")
	toTimestamp, err := ParseTimestamp("2018-02-01 00:00:00.000")

	if err != nil || !fromTimestamp.Time.Equal(rd.From.Time) || !toTimestamp.Time.Equal(rd.To.Time) {
		t.Error("Expected equal times")
	}
}

func TestGenerateTimeFrames(t *testing.T) {
	ts20180101 := parseTimestamp("2018-01-01 00:00:00.000", t)
	ts20180103 := parseTimestamp("2018-01-03 00:00:00.000", t)
	ts20181231 := parseTimestamp("2018-12-31 23:59:59.999", t)

	duration1d := parseDuration("24h", t)
	duration1h := parseDuration("1h", t)

	var testGenerateTimeFrames = func(from, to Timestamp, chunkSize Duration, expected int) {
		period := Period{
			TimeFrame: TimeFrame{from, to},
			ChunkSize: chunkSize,
		}

		l := len(GenerateTimeFrames(period))

		if l != expected {
			t.Errorf("Expected %v, got %v\n", expected, l)
		}
	}

	testGenerateTimeFrames(ts20180101, ts20180103, duration1d, 2)
	testGenerateTimeFrames(ts20180101, ts20180103, duration1h, 48)
	testGenerateTimeFrames(ts20180101, ts20181231, duration1d, 365)
}

func TestTimestampToString(t *testing.T) {
	ts := parseTimestamp("2018-01-01 00:00:00.000", t)
	if ts.String() != "2018-01-01 00:00:00.000" {
		t.Error("Expected '2018-01-01 00:00:00.000', got", ts.String())
	}
}

func TestLastTimestampForChunkSize(t *testing.T) {
	ts := parseTimestamp("2018-01-01 00:00:00.000", t)
	chunkSize := parseDuration("24h", t)
	tsLast := ts.LastTimestampForChunkSize(chunkSize)
	if tsLast.String() != "2018-01-01 23:59:59.999" {
		t.Error("Expected '2018-01-01 23:59:59.999', got", tsLast.String())
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