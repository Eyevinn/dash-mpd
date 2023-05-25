package mpd_test

import (
	"testing"
	"time"

	"github.com/Eyevinn/dash-mpd/mpd"
	"github.com/stretchr/testify/require"
)

func TestSeconds2DurPtr(t *testing.T) {
	cases := []struct {
		seconds   float64
		wantedDur mpd.Duration
	}{
		{
			seconds: 0, wantedDur: 0,
		},
		{
			seconds: 0.5, wantedDur: mpd.Duration(500 * time.Millisecond),
		},
		{
			seconds: 1.3, wantedDur: mpd.Duration(1300 * time.Millisecond),
		},
	}

	for _, c := range cases {
		p := mpd.Seconds2DurPtrFloat64(c.seconds)
		require.Equal(t, c.wantedDur, *p)
	}

}
