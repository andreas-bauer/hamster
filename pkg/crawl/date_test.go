package crawl

import "testing"
import "fmt"

func TestGenerateCrawlRange(t *testing.T) {
	d1 := NewDate(2018, 1, 1)
	d2 := NewDate(2018, 1, 31)

	l := len(GenerateCrawlRange(d1, d2))
	if l != 31 {
		t.Error("Expected 31, got ", l)
	}

	l = len(GenerateCrawlRange(d1, d1))
	if l != 1 {
		t.Error("Expected 1, got ", l)
	}
}

func TestCountDaysInCrawlRange(t *testing.T) {
	d1 := NewDate(2018, 1, 1)
	d2 := NewDate(2018, 1, 31)

	c := CountDaysInCrawlRange(d1, d2)
	if c != 31 {
		t.Error("Expected 31, got ", c)
	}

	c = CountDaysInCrawlRange(d1, d1)
	if c != 1 {
		t.Error("Expected 1, got ", c)
	}
}

func TestMin(t *testing.T) {
	d := NewDate(2018, 1, 1)
	fs := d.Min()
	if fs != "2018-01-01 00:00:00.000" {
		t.Error("Expected '2018-01-01 00:00:00.000', got ", fs)
	}
}

func TestLastTime(t *testing.T) {
	d := NewDate(2018, 1, 1)
	fs := d.Max()
	if fs != "2018-01-01 23:59:59.999" {
		t.Error("Expected '2018-01-01 23:59:59.999', got ", fs)
	}
}

func TestMarshalJSON(t *testing.T) {
	d := NewDate(2018, 1, 1)
	fmt.Println(d.MarshalJSON())
}

func TestUnmarshalJSON(t *testing.T) {
	date := Date{}
	date.UnmarshalJSON([]byte("\"2018-01-01\""))
	fmt.Print(date)
}
/*
func (date Date) MarshalJSON() ([]byte, error) {
	return []byte(date.time.Format(parseJSONDate)), nil
}

func (date *Date) UnmarshalJSON(data []byte) (err error) {
	time, err := time.Parse(formatStringDate, string(data))
	*date = Date{time}
	return err
}*/