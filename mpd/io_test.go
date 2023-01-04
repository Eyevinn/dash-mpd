package mpd_test

import (
	"testing"

	"github.com/Eyevinn/dash-mpd/mpd"
	"github.com/stretchr/testify/require"
)

func TestDateTime(t *testing.T) {
	cases := []struct {
		inTimeS  float64
		dateTime mpd.DateTime
	}{
		{0, "1970-01-01T00:00:00Z"},
		{0.5, "1970-01-01T00:00:00.5Z"},
		{946684800, "2000-01-01T00:00:00Z"},
	}
	for _, tc := range cases {
		got := mpd.ConvertToDateTime(tc.inTimeS)
		require.Equal(t, tc.dateTime, got)
	}
}

func TestDateTimeS(t *testing.T) {
	cases := []struct {
		inTimeS  int64
		dateTime mpd.DateTime
	}{
		{0, "1970-01-01T00:00:00Z"},
		{946684800, "2000-01-01T00:00:00Z"},
	}
	for _, tc := range cases {
		got := mpd.ConvertToDateTimeS(tc.inTimeS)
		require.Equal(t, tc.dateTime, got)
	}
}
