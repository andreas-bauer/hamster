package crawl

import (
	"encoding/json"
	"time"
)

type Duration struct {
	time.Duration
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
