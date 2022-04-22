package julianday

import (
	"errors"
	"math"
	"strconv"
	"time"
)

const secs_per_day = 86_400
const nsec_per_sec = 1_000_000_000
const nsec_per_day = nsec_per_sec * secs_per_day
const epoch_days = 2_440_587
const epoch_secs = secs_per_day / 2

func jd(t time.Time) (day, nsec int64) {
	sec := t.Unix()
	// guaranteed not to overflow
	day, sec = sec/secs_per_day+epoch_days, sec%secs_per_day+epoch_secs
	return day, sec*nsec_per_sec + int64(t.Nanosecond())
}

func Date(t time.Time) (day, nsec int64) {
	day, nsec = jd(t)
	switch {
	case nsec < 0:
		day -= 1
		nsec += nsec_per_day
	case nsec >= nsec_per_day:
		day += 1
		nsec -= nsec_per_day
	}
	return day, nsec
}

func Float(t time.Time) float64 {
	day, nsec := jd(t)
	// converting day and nsec to float64 is exact
	return float64(day) + float64(nsec)/nsec_per_day
}

func Format(t time.Time) string {
	var buf [32]byte
	return string(AppendFormat(buf[:0], t))
}

func AppendFormat(dst []byte, t time.Time) []byte {
	day, nsec := jd(t)
	if nsec >= nsec_per_day {
		day += 1
		nsec -= nsec_per_day
	}
	dst = strconv.AppendInt(dst, day, 10)
	pos := len(dst) - 1
	tmp := dst[pos]
	dst = strconv.AppendFloat(dst[:pos], math.Abs(float64(nsec))/nsec_per_day, 'f', -1, 64)
	dst[pos] = tmp
	return dst
}

func Time(day, nsec int64) time.Time {
	return time.Unix((day-epoch_days)*secs_per_day-epoch_secs, nsec)
}

func FloatTime(date float64) time.Time {
	day, frac := math.Modf(date)
	nsec := math.Floor(frac * nsec_per_day)
	return Time(int64(day), int64(nsec))
}

func Parse(s string) (time.Time, error) {
	dot := -1
	for i, b := range []byte(s) {
		if '0' <= b && b <= '9' {
			continue
		}
		if b == '.' && dot < 0 {
			dot = i
			continue
		}
		if (b == '+' || b == '-') && i == 0 {
			continue
		}
		return time.Time{}, errors.New("julianday: invalid syntax")
	}
	if len := len(s); len <= 1 && dot == len-1 {
		return time.Time{}, errors.New("julianday: invalid syntax")
	}

	var day, nsec int64
	if dot < 0 {
		day, _ = strconv.ParseInt(s, 10, 64)
	} else {
		if dot > 0 {
			day, _ = strconv.ParseInt(s[:dot], 10, 64)
		}
		frac, _ := strconv.ParseFloat(s[dot:], 64)
		nsec = int64(math.Floor(frac * nsec_per_day))
	}
	return Time(day, nsec), nil
}
