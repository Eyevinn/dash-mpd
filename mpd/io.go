package mpd

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/Eyevinn/dash-mpd/xml"
)

// countingWriter wraps an io.Writer and counts the bytes written.
type countingWriter struct {
	w io.Writer
	n int
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	cw.n += n
	return n, err
}

// ReadFromFile reads and unmarshals an MPD from a file.
func ReadFromFile(path string) (*MPD, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	mpd := MPD{}
	err = xml.Unmarshal(data, &mpd)
	if err != nil {
		return nil, err
	}
	mpd.SetParents()
	return &mpd, nil
}

// ReadFromString reads and unmarshals an MPD from a string
func ReadFromString(str string) (*MPD, error) {
	return MPDFromBytes([]byte(str))
}

// MPDFromBytes reads and unmarshals an MPD from a byte slice
func MPDFromBytes(data []byte) (*MPD, error) {
	mpd := MPD{}
	err := xml.Unmarshal(data, &mpd)
	if err != nil {
		return nil, err
	}
	mpd.SetParents()
	return &mpd, nil
}

// Write streams the XML encoding of m to w, with indentation according to indent
// and optionally prefixing an XML declaration header.  It returns the number of
// bytes written.  The encoder is pooled for efficiency; do not hold w open after
// Write returns.
func (m *MPD) Write(w io.Writer, indent string, withHeader bool) (int, error) {
	cw := &countingWriter{w: w}
	if withHeader {
		if _, err := io.WriteString(cw, xml.Header); err != nil {
			return cw.n, err
		}
	}
	enc := xml.AcquireEncoder(cw)
	defer xml.ReleaseEncoder(enc)
	enc.Indent("", indent)
	if err := enc.Encode(m); err != nil {
		return cw.n, err
	}
	if err := enc.Flush(); err != nil {
		return cw.n, err
	}
	return cw.n, nil
}

// WriteToString returns the XML encoding of m as a string, with indentation
// according to indent and optionally prefixing an XML declaration header.
func (m *MPD) WriteToString(indent string, withHeader bool) (string, error) {
	var sb strings.Builder
	sb.Grow(4096)
	if _, err := m.Write(&sb, indent, withHeader); err != nil {
		return "", err
	}
	return sb.String(), nil
}

const RFC3339MS string = "2006-01-02T15:04:05.999Z07:00"

// ConvertToDateTime converts a number of seconds to a UTC DateTime by cropping to ms precision.
func ConvertToDateTime(seconds float64) DateTime {
	s := int64(seconds)
	ns := int64((seconds - float64(s)) * 1_000_000_000)
	t := time.Unix(s, ns).UTC()

	return DateTime(t.Format(RFC3339MS))
}

// ConvertToDateTime converts an integral number of seconds to a UTC DateTime.
func ConvertToDateTimeS(seconds int64) DateTime {
	t := time.Unix(seconds, 0).UTC()
	return DateTime(t.Format(RFC3339MS))
}

// ConvertToDateTime converts an integral number of milliseconds to a UTC DateTime.
func ConvertToDateTimeMS(ms int64) DateTime {
	seconds := ms / 1000
	ns := (ms - 1000*seconds) * 1_000_000
	t := time.Unix(seconds, ns).UTC()
	return DateTime(t.Format(RFC3339MS))
}

// ConvertToSeconds converts a DateTime to a number of seconds.
func (dt DateTime) ConvertToSeconds() (float64, error) {
	t, err := time.Parse(RFC3339MS, string(dt))
	if err != nil {
		return 0, err
	}
	return float64(t.UnixNano()) / 1_000_000_000, nil
}
