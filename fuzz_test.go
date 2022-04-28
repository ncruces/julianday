//go:build go1.18

package julianday_test

import (
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/ncruces/julianday"
)

const secs_per_day = 86_400
const nsec_per_sec = 1_000_000_000
const nsec_per_day = nsec_per_sec * secs_per_day

func FuzzUnixDateTime(f *testing.F) {
	add := func(sec, nsec int64) { f.Add(sec, nsec) }

	add(0, 0)
	add(math.MinInt64, 0)
	add(math.MaxInt64, 0)
	add(math.MinInt64, 999_999_999)
	add(math.MaxInt64, 999_999_999)

	f.Fuzz(func(t *testing.T, sec, nsec int64) {
		tm := time.Unix(sec, nsec).UTC()
		if got := julianday.Time(julianday.Date(tm)); !tm.Equal(got) {
			t.Errorf("Time(Date(%v)) = %v", tm, got)
		}
	})
}

func FuzzUnixFormatParse(f *testing.F) {
	add := func(sec, nsec int64) { f.Add(sec, nsec) }

	add(0, 0)
	add(math.MinInt64, 0)
	add(math.MaxInt64, 0)
	add(math.MinInt64, 999_999_999)
	add(math.MaxInt64, 999_999_999)

	f.Fuzz(func(t *testing.T, sec, nsec int64) {
		tm := time.Unix(sec, nsec).UTC()
		got, err := julianday.Parse(julianday.Format(tm))
		if err != nil {
			t.Errorf("Parse(Format(%v)) = %v", tm, err)
		}
		if !tm.Equal(got) {
			t.Errorf("Parse(Format(%v)) = %v", tm, got)
		}
	})
}

func FuzzFloatTime(f *testing.F) {
	const minJD = -1e14
	const maxJD = +1e14
	add := func(jd float64) { f.Add(jd) }

	add(0)
	add(minJD)
	add(maxJD)

	f.Fuzz(func(t *testing.T, f float64) {
		if f < minJD || f > maxJD {
			t.SkipNow()
		}
		got := julianday.Float(julianday.FloatTime(f))
		if !nearlyEqual(f, got, 2, 2.0/nsec_per_day) { // nanosecond accuracy
			t.Errorf("Float(FloatTime(%g)) = %g", f, got)
		}
	})
}

func FuzzParseFloat(f *testing.F) {
	invalids := [...]string{"", ".", "+", "-", "..", "+.", "0+", "am", "10000000000000000000"}
	for _, s := range invalids {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		parsed, err := julianday.Parse(s)
		if err != nil && !parsed.IsZero() {
			t.Errorf("Parse(%q) = (%v, %v)", s, parsed, err)
		}
		if err != nil {
			t.SkipNow()
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			t.Errorf("ParseFloat(%q) = %v", s, err)
		}

		const minJD = -1e14
		const maxJD = +1e14
		if f < minJD || f > maxJD {
			t.SkipNow()
		}
		got := julianday.Float(parsed)
		if !nearlyEqual(f, got, 2, 2.0/nsec_per_day) { // nanosecond accuracy
			t.Errorf("Float(Parse(%q)) = %g", s, got)
		}
	})
}

func nearlyEqual(a, b float64, ulps int, abs float64) bool {
	if a == b {
		return true
	}
	if math.Abs(a-b) < abs {
		return true
	}
	for i := 0; i < ulps; i++ {
		a = math.Nextafter(a, b)
	}
	return a == b
}
