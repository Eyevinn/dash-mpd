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
		{0.333, "1970-01-01T00:00:00.333Z"},
		{0.666, "1970-01-01T00:00:00.666Z"},
		{-1, "1969-12-31T23:59:59Z"},
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
		{-1, "1969-12-31T23:59:59Z"},
		{946684800, "2000-01-01T00:00:00Z"},
	}
	for _, tc := range cases {
		got := mpd.ConvertToDateTimeS(tc.inTimeS)
		require.Equal(t, tc.dateTime, got)
	}
}

func TestDateTimeMS(t *testing.T) {
	cases := []struct {
		inTimeMS int64
		dateTime mpd.DateTime
	}{
		{0, "1970-01-01T00:00:00Z"},
		{-500, "1969-12-31T23:59:59.5Z"},
		{946684800120, "2000-01-01T00:00:00.12Z"},
	}
	for _, tc := range cases {
		got := mpd.ConvertToDateTimeMS(tc.inTimeMS)
		require.Equal(t, tc.dateTime, got)
	}
}

func TestConvertToSeconds(t *testing.T) {
	cases := []struct {
		desc     string
		dateTime mpd.DateTime
		wanted   float64
		err      string
	}{
		{
			desc:     "epoch start",
			dateTime: mpd.DateTime("1970-01-01T00:00:00Z"),
			wanted:   0,
			err:      "",
		},
		{
			desc:     "3.84s",
			dateTime: mpd.DateTime("1970-01-01T00:00:03.84Z"),
			wanted:   3.84,
			err:      "",
		},
		{
			desc:     "bad time",
			dateTime: mpd.DateTime("1970-01-01T00:00"),
			wanted:   0,
			err:      `parsing time "1970-01-01T00:00" as "2006-01-02T15:04:05.999Z07:00": cannot parse "" as ":"`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := tc.dateTime.ConvertToSeconds()
			if tc.err == "" {
				require.NoError(t, err)
				require.Equal(t, tc.wanted, got)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.err, err.Error())
			}
		})
	}
}
