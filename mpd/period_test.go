package mpd_test

import (
	"testing"

	"github.com/Eyevinn/dash-mpd/mpd"
	"github.com/stretchr/testify/require"
)

func TestPeriodStart(t *testing.T) {
	m := mpd.NewMPD()
	m.Type = mpd.Ptr("dynamic")
	m.AvailabilityStartTime = mpd.DateTime("1970-01-01T00:00:00Z")
	for i := 0; i < 3; i++ {
		p := mpd.NewPeriod()
		p.Start = mpd.Seconds2DurPtr(60 * i)
		m.Periods = append(m.Periods, p)
	}
	testCases := []struct {
		desc        string
		period      *mpd.PeriodType
		wantedStart float64
		err         string
	}{
		{"period 0", m.Periods[0], 0, ""},
		{"period 0", m.Periods[1], 60, ""},
		{"period 0", m.Periods[2], 120, ""},
		{"period 0", mpd.NewPeriod(), 0, "period not found in MPD"},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			gotStart, err := tc.period.AbsoluteStart(m)
			if tc.err != "" {
				require.EqualError(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantedStart, gotStart)
			}
		})
	}
}
