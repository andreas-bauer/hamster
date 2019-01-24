package crawl

import (
	"encoding/json"
	"time"
)

type Timestamp struct {
	time.Time
}

type TimeFrame struct {
	From Timestamp `json:"from"`
	To   Timestamp `json:"to"`
}

type Period struct {
	TimeFrame
	ChunkSize Duration `json:"chunkSize"`
}

type Duration struct {
	time.Duration
}

func GenerateTimeFrames(period Period) []TimeFrame {
	res := make([]TimeFrame, 0)

	for t := period.From.Time; t.Before(period.To.Time); t = t.Add(period.ChunkSize.Duration) {
		r := TimeFrame{}
		r.From.Time = t
		r.To.Time = t.Add(period.ChunkSize.Duration).Add(-1 * time.Millisecond)
		res = append(res, r)
	}
	return res
}

func ParseDuration(value string) (Duration, error) {
	t, err := time.ParseDuration(value)
	return Duration{t}, err
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	var err error
	d.Duration, err = time.ParseDuration(value)
	if err != nil {
		return err
	}
	return nil
}

const stringLayout = "2006-01-02 15:04:05.000"
const jsonLayout = "\"2006-01-02 15:04:05.000\""

func ParseTimestamp(value string) (Timestamp, error) {
	t, err := time.Parse(stringLayout, value)
	return Timestamp{t}, err
}

func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(ts.Time.Format(jsonLayout)), nil
}

func (ts *Timestamp) UnmarshalJSON(data []byte) (err error) {
	time, err := time.Parse(jsonLayout, string(data))
	ts.Time = time
	return err
}

func (ts *Timestamp) String() string {
	return ts.Time.Format(stringLayout)
}

func (ts Timestamp) LastTimestampForChunkSize(chunkSize Duration) Timestamp {
	return Timestamp{ts.Time.Add(chunkSize.Duration).Add(-1 * time.Millisecond)}
}
