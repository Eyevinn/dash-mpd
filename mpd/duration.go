package mpd

import (
	"math"
	"strconv"
	"time"

	"github.com/Eyevinn/dash-mpd/xml"
)

// ParseDurationError for parsing xs:Duration string
type ParseDurationError struct {
	msg string
}

func (p ParseDurationError) Error() string {
	return p.msg
}

func newParseDurationError(msg string) ParseDurationError {
	return ParseDurationError{msg: msg}
}

// This file is essentially a copy of https://github.com/zencoder/go-dash/blob/master/mpd/duration.go,
// adopted to work with the patched XML parser.
// Another change is that rendered times ending with 0S or 0M0S have that part removed.
// Copyright BrightCove Inc under Apache v2 License.

// Duration is an alias of time.Duration and has nano-second precision from Epoch start.
//
// XML marshaling methods need Duration to be included as a pointer in XML.
type Duration time.Duration

func (d *Duration) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if d == nil {
		return xml.Attr{}, nil
	}
	return xml.Attr{Name: name, Value: d.String()}, nil
}

func (d *Duration) UnmarshalXMLAttr(attr xml.Attr) error {
	dur, err := ParseDuration(attr.Value)
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

// String renders a Duration in XML Duration Data Type format.
//
// It handles negative durations, although they should not occur.
// The highest output unit is hours (H).
// There is never more than 3 decimals to the seconds.
func (d *Duration) String() string {
	// Largest time is 2540400h10m10.000s
	var buf [32]byte
	w := len(buf)

	u := uint64(*d)
	if u == 0 {
		return "PT0S"
	}
	neg := *d < 0
	if neg {
		u = -u
	}

	s := u / uint64(time.Second)
	ns := u - s*uint64(time.Second)
	ms := uint64(math.Round(float64(ns) * 1.0e-6))

	w--
	buf[w] = 'S' // End with Seconds

	switch {
	case s == 0 && ms == 0:
		// Time smaller than ms, return higher precision
		w, u = fmtFrac(buf[:w], u, 9)
		w = fmtInt(buf[:w], u)
	case s == 0:
		// Time smaller than 1s, return ms
		w, _ = fmtFrac(buf[:w], ms, 3)
		w--
		buf[w] = '0'
	default:
		// Time larger than 1s, return s and potentially ms
		u = 1000*s + ms
		w, u = fmtFrac(buf[:w], u, 3)
		// u is now integer seconds
		w = fmtInt(buf[:w], u%60)
		u /= 60
		if string(buf[w:]) == "0S" {
			w = len(buf) // Reset
		}

		// u is now integer minutes
		if u > 0 {
			w--
			buf[w] = 'M'
			w = fmtInt(buf[:w], u%60)
			u /= 60
			if string(buf[w:]) == "0M" {
				w = len(buf) // Reset
			}

			// u is now integer hours
			// Stop at hours because days can be different lengths.
			if u > 0 {
				w--
				buf[w] = 'H'
				w = fmtInt(buf[:w], u)
			}
		}
	}

	w--
	buf[w] = 'T'
	w--
	buf[w] = 'P'
	if neg {
		w--
		buf[w] = '-'
	}

	return string(buf[w:])
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros.  it omits the decimal
// point too when the fraction is 0.  It returns the index where the
// output bytes begin and the value v/10**prec.
func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for range prec {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return w, v
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

// ParseDuration parses an xsd Duration string and returns corresponding time.Duration.
func ParseDuration(str string) (time.Duration, error) {
	if len(str) < 3 {
		return 0, newParseDurationError("at least one number and designator are required")
	}
	if str[0] == '-' {
		return 0, newParseDurationError("duration cannot be negative")
	}
	if str[0] != 'P' {
		return 0, newParseDurationError("duration must be in the format: P[nD][T[nH][nM][nS]]")
	}

	var (
		total    time.Duration
		i        = 1 // cursor, positioned just after 'P'
		afterT   = false
		seenAny  = false
		lastUnit byte
	)

	// rank enforces D < H < M < S ordering.
	rank := func(c byte) int {
		switch c {
		case 'D':
			return 1
		case 'H':
			return 2
		case 'M':
			return 3
		case 'S':
			return 4
		}
		return 0
	}

	for i < len(str) {
		c := str[i]
		if c == 'T' {
			if afterT {
				return 0, newParseDurationError("duration must be in the format: P[nD][T[nH][nM][nS]]")
			}
			afterT = true
			i++
			continue
		}
		if c == '-' {
			return 0, newParseDurationError("duration cannot be negative")
		}
		// Must now read digits (and maybe '.', only if looking for S).
		start := i
		sawDot := false
		for i < len(str) {
			b := str[i]
			if b >= '0' && b <= '9' {
				i++
				continue
			}
			if b == '.' && !sawDot {
				sawDot = true
				i++
				continue
			}
			break
		}
		if i == start {
			return 0, newParseDurationError("duration must be in the format: P[nD][T[nH][nM][nS]]")
		}
		if i == len(str) {
			return 0, newParseDurationError("duration must be in the format: P[nD][T[nH][nM][nS]]")
		}
		unit := str[i]
		i++
		if rank(unit) == 0 {
			return 0, newParseDurationError("duration must be in the format: P[nD][T[nH][nM][nS]]")
		}
		if rank(unit) <= rank(lastUnit) {
			return 0, newParseDurationError("duration must be in the format: P[nD][T[nH][nM][nS]]")
		}
		if unit != 'D' && !afterT {
			return 0, newParseDurationError("duration must be in the format: P[nD][T[nH][nM][nS]]")
		}
		if sawDot && unit != 'S' {
			return 0, newParseDurationError("duration must be in the format: P[nD][T[nH][nM][nS]]")
		}
		lastUnit = unit
		seenAny = true

		numStr := str[start : i-1] // digits (and maybe '.') without the unit letter

		switch unit {
		case 'D':
			n, err := strconv.ParseUint(numStr, 10, 64)
			if err != nil {
				return 0, newParseDurationError("error parsing Days")
			}
			total += time.Duration(n) * 24 * time.Hour
		case 'H':
			n, err := strconv.ParseUint(numStr, 10, 64)
			if err != nil {
				return 0, newParseDurationError("error parsing Hours")
			}
			total += time.Duration(n) * time.Hour
		case 'M':
			n, err := strconv.ParseUint(numStr, 10, 64)
			if err != nil {
				return 0, newParseDurationError("error parsing Minutes")
			}
			total += time.Duration(n) * time.Minute
		case 'S':
			f, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, newParseDurationError("error parsing Seconds")
			}
			total += time.Duration(f * float64(time.Second))
		}
	}
	if !seenAny {
		return 0, newParseDurationError("at least one number and designator are required")
	}
	return total, nil
}

// Seconds returns the duration as a floating point number of seconds.
func (d *Duration) Seconds() float64 {
	return float64(*d) / float64(time.Second)
}
