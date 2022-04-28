package julianday

import (
	"strconv"
	"testing"
	"time"
)

var reference = time.Date(1969, 7, 20, 22, 17, 40, 0, time.UTC)

var years = [...]struct {
	year int
	want int64
}{
	{-5000, -105152},
	{-4000, +260090},
	{-0001, 1721426 - 366 - 365},
	{+0000, 1721426 - 366},
	{+0001, 1721426},
	{+1000, 2086303},
	{+2000, 2451545},
}

func TestDate(t *testing.T) {
	gotDay, gotNsec := Date(reference)
	if gotDay != 2440423 {
		t.Errorf("Date() gotDay = %d, want %d", gotDay, 2440423)
	}
	if gotNsec != 37060000000000 {
		t.Errorf("Date() gotNsec = %d, want %d", gotNsec, 37060000000000)
	}
}

func TestDate_midday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 12, 0, 0, 0, time.UTC)
		gotDay, gotNsec := Date(date)
		if gotDay != tt.want {
			t.Errorf("Date() gotDay = %d, want %d", gotDay, tt.want)
		}
		if gotNsec != 0 {
			t.Errorf("Date() gotNsec = %d, want %d", gotNsec, 0)
		}
	}
}

func TestDate_afterMidday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 12, 0, 0, 1, time.UTC)
		gotDay, gotNsec := Date(date)
		if gotDay != tt.want {
			t.Errorf("Date() gotDay = %d, want %d", gotDay, tt.want)
		}
		if gotNsec != 1 {
			t.Errorf("Date() gotNsec = %d, want %d", gotNsec, 1)
		}
	}
}

func TestDate_beforeMidday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 11, 59, 59, nsec_per_sec-1, time.UTC)
		gotDay, gotNsec := Date(date)
		if gotDay != tt.want-1 {
			t.Errorf("Date() gotDay = %d, want %d", gotDay, tt.want-1)
		}
		if gotNsec != nsec_per_day-1 {
			t.Errorf("Date() gotNsec = %d, want %d", gotNsec, nsec_per_day-1)
		}
	}
}

func TestFloat(t *testing.T) {
	got := Float(reference)
	if got != 2440423.428935185185185 {
		t.Errorf("Float() got = %f, want %f", got, 2440423.428935185185185)
	}
}

func TestFloat_sqlite(t *testing.T) {
	// https://www.sqlite.org/lang_datefunc.html
	date, _ := time.Parse(time.RFC3339Nano, "2013-10-07T08:23:19.120Z")
	got := Float(date)
	if got != 2456572.849526852 {
		t.Errorf("Float() got = %g, want %g", got, 2456572.849526852)
	}
}

func TestFloat_midday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 12, 0, 0, 0, time.UTC)
		got := Float(date)
		if got != float64(tt.want) {
			t.Errorf("Float() got = %f, want %d", got, tt.want)
		}
	}
}

func TestFloat_afterMidday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 12, 0, 0, 1, time.UTC)
		got := Float(date)
		if got != float64(tt.want) {
			t.Errorf("Float() got = %f, want %d", got, tt.want)
		}
	}
}

func TestFloat_beforeMidday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 11, 59, 59, nsec_per_sec-1, time.UTC)
		got := Float(date)
		if got != float64(tt.want) {
			t.Errorf("Float() got = %f, want %d", got, tt.want)
		}
	}
}

func TestFormat(t *testing.T) {
	got := Format(reference)
	if got != "2440423.428935185185185" {
		t.Errorf("Format() got = %s, want %s", got, "2440423.428935185185185")
	}
}

func TestFormat_midday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 12, 0, 0, 0, time.UTC)
		want := strconv.FormatInt(tt.want, 10)
		got := Format(date)
		if got != want {
			t.Errorf("Format() got = %s, want %s", got, want)
		}
	}
}

func TestFormat_afterMidday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 12, 0, 0, 1, time.UTC)
		other, err := strconv.ParseFloat(Format(date), 64)
		if err != nil {
			t.Errorf("Format() got = %v", err)
		}
		if other != float64(tt.want) {
			t.Errorf("Format() got = %f, want %d", other, tt.want)
		}
	}
}

func TestFormat_beforeMidday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 11, 59, 59, nsec_per_sec-1, time.UTC)
		other, err := strconv.ParseFloat(Format(date), 64)
		if err != nil {
			t.Errorf("Format() got = %v", err)
		}
		if other != float64(tt.want) {
			t.Errorf("Format() got = %f, want %d", other, tt.want)
		}
	}
}
func TestTime(t *testing.T) {
	got := Time(2440423, 37060000000000)
	if !got.Equal(reference) {
		t.Errorf("Time() got = %v, want %v", got, reference)
	}
}

func TestTime_midday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 12, 0, 0, 0, time.UTC)
		day, nsec := Date(date)
		got := Time(day, nsec)
		if !got.Equal(date) {
			t.Errorf("Time() got = %v, want %v", got, date)
		}
	}
}

func TestTime_afterMidday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 12, 0, 0, 1, time.UTC)
		day, nsec := Date(date)
		got := Time(day, nsec)
		if !got.Equal(date) {
			t.Errorf("Time() got = %v, want %v", got, date)
		}
	}
}

func TestTime_beforeMidday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 11, 59, 59, nsec_per_sec-1, time.UTC)
		day, nsec := Date(date)
		got := Time(day, nsec)
		if !got.Equal(date) {
			t.Errorf("Time() got = %v, want %v", got, date)
		}
	}
}

func TestFloatTime(t *testing.T) {
	got := FloatTime(2440423.428935185185185)
	if got = got.Round(time.Millisecond); !got.Equal(reference) {
		t.Errorf("FloatTime() got = %v, want %v", got, reference)
	}
}

func TestParse(t *testing.T) {
	got, err := Parse("2440423.428935185185185")
	if err != nil {
		t.Errorf("Parse() got = %v", err)
	}
	if !got.Equal(reference) {
		t.Errorf("Parse() got = %v, want %v", got, reference)
	}
}

func TestParse_date(t *testing.T) {
	got, err := Parse("2440423")
	if err != nil {
		t.Errorf("Parse() got = %v", err)
	}
	goty, gotm, gotd := got.Date()
	refy, refm, refd := reference.Date()
	if goty != refy || gotm != refm || gotd != refd {
		t.Errorf("Parse() got = %v, want %v", got, reference)
	}
}

func TestParse_clock(t *testing.T) {
	got, err := Parse(".428935185185185")
	if err != nil {
		t.Errorf("Parse() got = %v", err)
	}
	goth, gotm, gots := got.Clock()
	refh, refm, refs := reference.Clock()
	if goth != refh || gotm != refm || gots != refs {
		t.Errorf("Parse() got = %v, want %v", got, reference)
	}
}

func TestParse_invalid(t *testing.T) {
	invalids := [...]string{"", ".", "+", "-", "..", "+.", "0+", "am", "10000000000000000000"}
	for _, s := range invalids {
		if got, err := Parse(s); err == nil {
			t.Errorf("Parse() got = (%v, %v)", got, err)
		}
	}
}

func TestParse_midday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 12, 0, 0, 0, time.UTC)
		str := Format(date)
		got, err := Parse(str)
		if err != nil {
			t.Errorf("Parse() got = %v", err)
		}
		if !got.Equal(date) {
			t.Errorf("Parse() got = %v, want %v", got, date)
		}
	}
}

func TestParse_afterMidday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 12, 0, 0, 1, time.UTC)
		str := Format(date)
		got, err := Parse(str)
		if err != nil {
			t.Errorf("Parse() got = %v", err)
		}
		if !got.Equal(date) {
			t.Errorf("Parse() got = %v, want %v", got, date)
		}
	}
}

func TestParse_beforeMidday(t *testing.T) {
	for _, tt := range years {
		date := time.Date(tt.year, 1, 1, 11, 59, 59, nsec_per_sec-1, time.UTC)
		str := Format(date)
		got, err := Parse(str)
		if err != nil {
			t.Errorf("Parse() got = %v", err)
		}
		if !got.Equal(date) {
			t.Errorf("Parse() got = %v, want %v", got, date)
		}
	}
}
