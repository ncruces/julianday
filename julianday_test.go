package julianday_test

import (
	"testing"
	"time"

	"github.com/ncruces/julianday"
)

var reference = time.Date(1969, 7, 20, 22, 17, 40, 0, time.UTC)

func TestDate(t *testing.T) {
	gotDay, gotNsec := julianday.Date(reference)
	if gotDay != 2440423 {
		t.Errorf("Date() gotDay = %v, want %v", gotDay, 2440423)
	}
	if gotNsec != 37060000000000 {
		t.Errorf("Date() gotNsec = %v, want %v", gotNsec, 37060000000000)
	}
}

func TestFloat(t *testing.T) {
	got := julianday.Float(reference)
	if got != 2440423.4289351851851852 {
		t.Errorf("Float() got = %f, want %f", got, 2440423.4289351851851852)
	}
}

func TestFormat(t *testing.T) {
	got := julianday.Format(reference)
	if got != "2440423.4289351851851852" {
		t.Errorf("Format() got = %s, want %s", got, "2440423.4289351851851852")
	}
}

func TestTime(t *testing.T) {
	got := julianday.Time(2440423, 37060000000000)
	if !got.Equal(reference) {
		t.Errorf("Time() got = %v, want %v", got, reference)
	}
}

func TestFloatTime(t *testing.T) {
	got := julianday.FloatTime(2440423.4289351851851852)
	if got = got.Round(time.Millisecond); !got.Equal(reference) {
		t.Errorf("FloatTime() got = %v, want %v", got, reference)
	}
}

func TestParse(t *testing.T) {
	got, err := julianday.Parse("2440423.4289351851851852")
	if err != nil {
		t.Errorf("Parse() got = %v", err)
	}
	if !got.Equal(reference) {
		t.Errorf("Parse() got = %v, want %v", got, reference)
	}
}
