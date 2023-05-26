package mpd_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Eyevinn/dash-mpd/mpd"
	"github.com/stretchr/testify/require"
)

func TestAbsolutePeriodStart(t *testing.T) {
	m := mpd.NewMPD(mpd.DYNAMIC_TYPE)
	m.AvailabilityStartTime = mpd.DateTime("1970-01-01T00:00:00Z")
	for i := 0; i < 3; i++ {
		p := mpd.NewPeriod()
		p.Start = mpd.Seconds2DurPtr(60 * i)
		m.AppendPeriod(p)
	}
	m.SetParents()
	testCases := []struct {
		desc        string
		period      *mpd.Period
		wantedStart float64
		err         string
	}{
		{"period 0", m.Periods[0], 0, ""},
		{"period 0", m.Periods[1], 60, ""},
		{"period 0", m.Periods[2], 120, ""},
		{"period 0", mpd.NewPeriod(), 0, mpd.ErrParentNotSet.Error()},
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

func TestPeriodStart(t *testing.T) {
	testCases := []struct {
		desc         string
		mpdType      string
		mupPresent   bool
		data         []mpdData
		wantedStarts []mpd.Duration
	}{
		{
			desc:         "dynamic, single-period",
			mpdType:      mpd.DYNAMIC_TYPE,
			mupPresent:   true,
			data:         []mpdData{{mpd.Seconds2DurPtr(12), nil}},
			wantedStarts: []mpd.Duration{mpd.Duration(12 * time.Second)},
		},
		{
			desc:       "dynamic, single-period",
			mpdType:    mpd.DYNAMIC_TYPE,
			mupPresent: true,
			data: []mpdData{
				{mpd.Seconds2DurPtr(12), mpd.Seconds2DurPtr(60)},
				{nil, nil},
			},
			wantedStarts: []mpd.Duration{
				mpd.Duration(12 * time.Second),
				mpd.Duration(72 * time.Second),
			},
		},
		{
			desc:       "static, single-period without start",
			mpdType:    mpd.STATIC_TYPE,
			mupPresent: false,
			data: []mpdData{
				{nil, nil},
			},
			wantedStarts: []mpd.Duration{
				mpd.Duration(0 * time.Second),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			m := mpd.NewMPD(tc.mpdType)
			m.AvailabilityStartTime = mpd.DateTime("1970-01-01T00:00:00Z")
			for _, d := range tc.data {
				p := mpd.NewPeriod()
				p.Start = d.start
				p.Duration = d.dur
				m.AppendPeriod(p)
			}
			m.SetParents()
			for i, p := range m.Periods {
				gotStart, err := p.GetStart()
				require.NoError(t, err)
				require.Equal(t, tc.wantedStarts[i], gotStart)
			}
		})
	}
}

type mpdData struct {
	start *mpd.Duration
	dur   *mpd.Duration
}

func TestMpdStartErrors(t *testing.T) {
	m := mpd.NewMPD(mpd.DYNAMIC_TYPE)
	p := mpd.NewPeriod()
	m.Periods = append(m.Periods, p)
	_, err := p.GetStart()
	require.EqualError(t, err, mpd.ErrParentNotSet.Error())

	m = mpd.NewMPD(mpd.DYNAMIC_TYPE)
	p = mpd.NewPeriod()
	m.AppendPeriod(p)
	m.Periods = nil
	_, err = p.GetStart()
	require.EqualError(t, err, mpd.ErrPeriodNotFound.Error())

	m = mpd.NewMPD(mpd.DYNAMIC_TYPE)
	p = mpd.NewPeriod()
	m.AppendPeriod(p)
	p = mpd.NewPeriod()
	m.AppendPeriod(p)
	_, err = p.GetStart()
	require.EqualError(t, err, mpd.ErrUnknownPeriodStart.Error())

	m = mpd.NewMPD(mpd.DYNAMIC_TYPE)
	p = mpd.NewPeriod()
	p.Duration = mpd.Seconds2DurPtr(30)
	m.AppendPeriod(p)
	p = mpd.NewPeriod()
	m.AppendPeriod(p)
	_, err = p.GetStart()
	require.EqualError(t, err, "cannot get start of previous period: period start cannot be derived")
}

func TestPeriodType(t *testing.T) {

	testCases := []struct {
		desc        string
		mpdType     string
		mupPresent  bool
		data        []mpdData
		wantedTypes []mpd.PeriodType
	}{
		{
			desc:        "dynamic, single-period",
			mpdType:     mpd.DYNAMIC_TYPE,
			mupPresent:  false,
			data:        []mpdData{{mpd.Seconds2DurPtr(0), nil}},
			wantedTypes: []mpd.PeriodType{mpd.PTRegular},
		},
		{
			desc:        "static, single-period",
			mpdType:     mpd.STATIC_TYPE,
			mupPresent:  false,
			data:        []mpdData{{mpd.Seconds2DurPtr(0), nil}},
			wantedTypes: []mpd.PeriodType{mpd.PTRegular},
		},
		{
			desc:       "dynamic, first with dur, next without start and dur",
			mpdType:    mpd.DYNAMIC_TYPE,
			mupPresent: true,
			data: []mpdData{
				{mpd.Seconds2DurPtr(0), mpd.Seconds2DurPtr(20)},
				{nil, nil},
			},
			wantedTypes: []mpd.PeriodType{mpd.PTEarlyTerminated, mpd.PTRegularOrEarlyTerminated},
		},
		{
			desc:       "dynamic, first with dur, next without start and dur",
			mpdType:    mpd.DYNAMIC_TYPE,
			mupPresent: true,
			data: []mpdData{
				{mpd.Seconds2DurPtr(0), mpd.Seconds2DurPtr(20)},
				{nil, nil},
			},
			wantedTypes: []mpd.PeriodType{mpd.PTEarlyTerminated, mpd.PTRegularOrEarlyTerminated},
		},
		{
			desc:       "dynamic, early available first segment",
			mpdType:    mpd.DYNAMIC_TYPE,
			mupPresent: true,
			data: []mpdData{
				{nil, nil},
			},
			wantedTypes: []mpd.PeriodType{mpd.PTEarlyAvailable},
		},
		{
			desc:       "dynamic, early available second segment",
			mpdType:    mpd.DYNAMIC_TYPE,
			mupPresent: true,
			data: []mpdData{
				{mpd.Seconds2DurPtr(0), nil},
				{nil, nil},
			},
			wantedTypes: []mpd.PeriodType{mpd.PTRegular, mpd.PTEarlyAvailable},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			m := mpd.NewMPD(tc.mpdType)
			if tc.mupPresent {
				m.MinimumUpdatePeriod = mpd.Seconds2DurPtr(2)
			}
			for _, d := range tc.data {
				p := mpd.NewPeriod()
				p.Start = d.start
				p.Duration = d.dur
				m.AppendPeriod(p)
			}
			m.SetParents()
			for i, p := range m.Periods {
				pType, err := p.GetType()
				require.NoError(t, err)
				require.Equal(t, tc.wantedTypes[i], pType, fmt.Sprintf("period %d", i))
			}
		})
	}
}

func TestMpdTypeErrors(t *testing.T) {
	m := mpd.NewMPD(mpd.DYNAMIC_TYPE)
	p := mpd.NewPeriod()
	m.Periods = append(m.Periods, p)
	_, err := p.GetType()
	require.EqualError(t, err, mpd.ErrParentNotSet.Error())

	m = mpd.NewMPD(mpd.DYNAMIC_TYPE)
	p = mpd.NewPeriod()
	m.AppendPeriod(p)
	m.Periods = nil
	_, err = p.GetType()
	require.EqualError(t, err, mpd.ErrPeriodNotFound.Error())
}

func TestGetDuration(t *testing.T) {

	testCases := []struct {
		desc         string
		mpdType      string
		mediaPresDur int
		data         []mpdData
		wantedDurs   []mpd.Duration
		wantedErrs   []error
	}{
		{
			desc:         "static with duration",
			mpdType:      mpd.STATIC_TYPE,
			mediaPresDur: 0,
			data:         []mpdData{{nil, mpd.Seconds2DurPtr(60)}},
			wantedDurs:   []mpd.Duration{mpd.Duration(60 * time.Second)},
			wantedErrs:   nil,
		},
		{
			desc:         "static without mediaPresendationDuration",
			mpdType:      mpd.STATIC_TYPE,
			mediaPresDur: 0,
			data:         []mpdData{{mpd.Seconds2DurPtr(0), nil}},
			wantedDurs:   nil,
			wantedErrs:   []error{mpd.ErrNoMediaPresentationDuration},
		},
		{
			desc:         "static not-last without next start",
			mpdType:      mpd.STATIC_TYPE,
			mediaPresDur: 120,
			data: []mpdData{
				{mpd.Seconds2DurPtr(0), nil},
				{nil, nil},
			},
			wantedDurs: nil,
			wantedErrs: []error{
				mpd.ErrUnknownPeriodDur,
				mpd.ErrUnknownPeriodDur,
			},
		},
		{
			desc:         "static with next start",
			mpdType:      mpd.STATIC_TYPE,
			mediaPresDur: 120,
			data: []mpdData{
				{mpd.Seconds2DurPtr(0), nil},
				{mpd.Seconds2DurPtr(90), nil},
			},
			wantedDurs: []mpd.Duration{
				mpd.Duration(90 * time.Second),
				mpd.Duration(30 * time.Second),
			},
			wantedErrs: nil,
		},
		{
			desc:         "dynamic, no start",
			mpdType:      mpd.DYNAMIC_TYPE,
			mediaPresDur: 0,
			data:         []mpdData{{nil, nil}},
			wantedDurs:   nil,
			wantedErrs:   []error{mpd.ErrUnknownPeriodDur},
		},
		{
			desc:         "dynamic, single-period",
			mpdType:      mpd.DYNAMIC_TYPE,
			mediaPresDur: 0,
			data:         []mpdData{{mpd.Seconds2DurPtr(0), nil}},
			wantedDurs:   nil,
			wantedErrs:   []error{mpd.ErrUnknownPeriodDur},
		},
		{
			desc:         "dynamic, without next start",
			mpdType:      mpd.DYNAMIC_TYPE,
			mediaPresDur: 0,
			data: []mpdData{
				{mpd.Seconds2DurPtr(0), nil},
				{nil, nil},
			},
			wantedDurs: nil,
			wantedErrs: []error{
				mpd.ErrUnknownPeriodDur,
				mpd.ErrUnknownPeriodDur,
			},
		},
		{
			desc:         "dynamic with next start",
			mpdType:      mpd.DYNAMIC_TYPE,
			mediaPresDur: 0,
			data: []mpdData{
				{mpd.Seconds2DurPtr(0), nil},
				{mpd.Seconds2DurPtr(90), nil},
			},
			wantedDurs: []mpd.Duration{
				mpd.Duration(90 * time.Second),
				0,
			},
			wantedErrs: []error{nil, mpd.ErrUnknownPeriodDur},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			m := mpd.NewMPD(tc.mpdType)
			if tc.mediaPresDur > 0 {
				m.MediaPresentationDuration = mpd.Seconds2DurPtr(tc.mediaPresDur)
			}
			for _, d := range tc.data {
				p := mpd.NewPeriod()
				p.Start = d.start
				p.Duration = d.dur
				m.AppendPeriod(p)
			}
			m.SetParents()
			for i, p := range m.Periods {
				dur, err := p.GetDuration()
				switch {
				case len(tc.wantedErrs) > 0 && tc.wantedErrs[i] != nil:
					require.EqualError(t, err, tc.wantedErrs[i].Error())
				case len(tc.wantedDurs) > 0:
					require.Equal(t, tc.wantedDurs[i], dur, fmt.Sprintf("period %d", i))
				default:
					t.Fail()
				}
			}
		})
	}
}
