//go:build go1.18

package julianday

import (
	"math"
	"strconv"
	"testing"
	"time"
)

func FuzzUnixDateTime(f *testing.F) {
	add := func(sec, nsec int64) { f.Add(sec, nsec) }

	add(0, 0)
	add(math.MinInt64, 0)
	add(math.MaxInt64, 0)
	add(math.MinInt64/2, 0)
	add(math.MaxInt64/2, 0)
	add(math.MinInt64, 999_999_999)
	add(math.MaxInt64, 999_999_999)

	f.Fuzz(func(t *testing.T, sec, nsec int64) {
		tm := time.Unix(sec, nsec).UTC()
		if got := Time(Date(tm)); !tm.Equal(got) {
			t.Errorf("Time(Date(%v)) = %v", tm, got)
		}
	})
}

func FuzzUnixFormatParse(f *testing.F) {
	add := func(sec, nsec int64) { f.Add(sec, nsec) }

	add(0, 0)
	add(math.MinInt64, 0)
	add(math.MaxInt64, 0)
	add(math.MinInt64/2, 0)
	add(math.MaxInt64/2, 0)
	add(math.MinInt64, 999_999_999)
	add(math.MaxInt64, 999_999_999)

	f.Fuzz(func(t *testing.T, sec, nsec int64) {
		tm := time.Unix(sec, nsec).UTC()
		got, err := Parse(Format(tm))
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
	add(2440423.428935185185185)
	add(2456572.849526852)

	f.Fuzz(func(t *testing.T, f float64) {
		if f < minJD || f > maxJD {
			t.SkipNow()
		}
		got := Float(FloatTime(f))
		if !nearlyEqual(f, got, 2, 2.0/nsec_per_day) { // nanosecond accuracy
			t.Errorf("Float(FloatTime(%g)) = %g", f, got)
		}
	})
}

func FuzzParseFloat(f *testing.F) {
	seed := [...]string{
		"2440423.428935185185185",
		"2440423", ".428935185185185",
		"", ".", "+", "-", "..", "+.", "0+", "am", "10000000000000000000",
	}
	for _, s := range seed {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		parsed, err := Parse(s)
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
		got := Float(parsed)
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
