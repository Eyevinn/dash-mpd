package mpd

import (
	"math"
	"regexp"
	"strconv"
	"strings"
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

const (
	rStart   = "^P"          // Must start with a 'P'
	rDays    = "(\\d+D)?"    // We only allow Days for durations, not Months or Years
	rTime    = "(?:T"        // If there's any 'time' units then they must be preceded by a 'T'
	rHours   = "(\\d+H)?"    // Hours
	rMinutes = "(\\d+M)?"    // Minutes
	rSeconds = "([\\d.]+S)?" // Seconds (Potentially decimal)
	rEnd     = ")?$"         // end of regex must close "T" capture group
)

var xmlDurationRegex = regexp.MustCompile(rStart + rDays + rTime + rHours + rMinutes + rSeconds + rEnd)

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

	if neg {
		w--
		buf[w] = '-'
	}

	return "PT" + string(buf[w:])
}

// fmtFrac formats the fraction of v/10**prec (e.g., ".12345") into the
// tail of buf, omitting trailing zeros.  it omits the decimal
// point too when the fraction is 0.  It returns the index where the
// output bytes begin and the value v/10**prec.
func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for i := 0; i < prec; i++ {
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

	if strings.Contains(str, "-") {
		return 0, newParseDurationError("duration cannot be negative")
	}

	// Check that only the parts we expect exist and that everything's in the correct order
	if !xmlDurationRegex.Match([]byte(str)) {
		return 0, newParseDurationError("duration must be in the format: P[nD][T[nH][nM][nS]]")
	}

	var parts = xmlDurationRegex.FindStringSubmatch(str)
	var total time.Duration

	if parts[1] != "" {
		days, err := strconv.Atoi(strings.TrimRight(parts[1], "D"))
		if err != nil {
			return 0, newParseDurationError("error parsing Days")
		}
		total += time.Duration(days) * time.Hour * 24
	}

	if parts[2] != "" {
		hours, err := strconv.Atoi(strings.TrimRight(parts[2], "H"))
		if err != nil {
			return 0, newParseDurationError("error parsing Hours")
		}
		total += time.Duration(hours) * time.Hour
	}

	if parts[3] != "" {
		mins, err := strconv.Atoi(strings.TrimRight(parts[3], "M"))
		if err != nil {
			return 0, newParseDurationError("error parsing Minutes")
		}
		total += time.Duration(mins) * time.Minute
	}

	if parts[4] != "" {
		secs, err := strconv.ParseFloat(strings.TrimRight(parts[4], "S"), 64)
		if err != nil {
			return 0, newParseDurationError("error parsing Seconds")
		}
		total += time.Duration(secs * float64(time.Second))
	}

	return total, nil
}

// Seconds returns the duration as a floating point number of seconds.
func (d *Duration) Seconds() float64 {
	return float64(*d) / float64(time.Second)
}
