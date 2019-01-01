package crawl

import (
	"time"
)

const formatStringDate = "2006-01-02"
const parseJSONDate = "\"2006-01-02\""

type Date struct {
	time time.Time
}

func NewDate(year, month, day int) Date {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return Date{t}
}

func GenerateCrawlRange(start, end Date) []Date {
	res := make([]Date, 0)
	for t := start.time; !t.After(end.time); t = t.AddDate(0, 0, 1) {
		res = append(res, Date{t})
	}
	return res
}

func (date Date) MarshalJSON() ([]byte, error) {
	return []byte(date.time.Format(parseJSONDate)), nil
}

func (date *Date) UnmarshalJSON(data []byte) (err error) {
	time, err := time.Parse(formatStringDate, string(data))
	*date = Date{time}
	return err
}

func CountDaysInCrawlRange(start, end Date) int {
	return int(end.time.Sub(start.time).Hours()/24) + 1
}

func (date Date) Min() string {
	return date.time.Format("2006-01-02 15:04:05.000")
}

func (date Date) Max() string {
	return date.time.AddDate(0, 0, 1).Add(-1 * time.Millisecond).Format("2006-01-02 15:04:05.000")
}
