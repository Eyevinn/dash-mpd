package mpd

import (
	"io"
	"os"
	"time"

	"github.com/Eyevinn/dash-mpd/xml"
)

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
	return &mpd, nil
}

// ReadFromString reads and unmarshals an MPD from a string
func ReadFromString(str string) (*MPD, error) {
	mpd := MPD{}
	err := xml.Unmarshal([]byte(str), &mpd)
	if err != nil {
		return nil, err
	}
	return &mpd, nil
}

// Write marshals and writes an MPD.
func (m *MPD) Write(w io.Writer) (int, error) {
	data, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		return 0, err
	}
	return w.Write(data)
}

// ConvertToDateTime converts a number of seconds to a UTC DateTime.
func ConvertToDateTime(seconds float64) DateTime {
	s := int64(seconds)
	ns := int64((seconds - float64(s)) * 1_000_000_000)
	t := time.Unix(s, ns).UTC()
	return DateTime(t.Format(time.RFC3339Nano))
}

// ConvertToDateTime converts an intergral number of seconds to a UTC DateTime.
func ConvertToDateTimeS(seconds int64) DateTime {
	t := time.Unix(seconds, 0).UTC()
	return DateTime(t.Format(time.RFC3339Nano))
}
