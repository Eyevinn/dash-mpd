package mpd_test

import (
	"strings"
	"testing"

	m "github.com/Eyevinn/dash-mpd/mpd"
	"github.com/stretchr/testify/require"
)

// buildAltReplaceMPD builds a minimal dynamic MPD with a single Alternative-MPD Replace event.
// The earliestResolutionTimeOffset is 5400000 (60s * 90000), the value that previously
// triggered scientific-notation rendering.
func buildAltReplaceMPD(clip *bool) *m.MPD {
	mpd := m.NewMPD(m.DYNAMIC_TYPE)
	p := m.NewPeriod()
	p.EventStreams = []*m.EventStreamType{{
		SchemeIdUri: "urn:mpeg:dash:event:alternativeMPD:replace:2025",
		Timescale:   m.Ptr(uint32(90000)),
		Events: []*m.EventType{{
			PresentationTime: 2700000,
			Id:               m.Ptr(uint64(1)),
			ReplacePresentation: &m.AlternativeMPDReplaceEventType{
				Clip: clip,
				AlternativeMPDEventType: &m.AlternativeMPDEventType{
					Uri:                          "http://example.com/ad.mpd",
					EarliestResolutionTimeOffset: 5400000,
					MaxDuration:                  1350000,
				},
			},
		}},
	}}
	mpd.AppendPeriod(p)
	return mpd
}

// TestAlternativeMPDReplaceClipRendering verifies that @clip is a *bool so that clip="false"
// can be rendered. With a plain bool+omitempty, clip="false" was silently dropped and a
// conforming client then applied the schema default (true) — the opposite of the intent.
func TestAlternativeMPDReplaceClipRendering(t *testing.T) {
	out, err := buildAltReplaceMPD(m.Ptr(false)).WriteToString("  ", false)
	require.NoError(t, err)
	require.Contains(t, out, `clip="false"`)

	out, err = buildAltReplaceMPD(m.Ptr(true)).WriteToString("  ", false)
	require.NoError(t, err)
	require.Contains(t, out, `clip="true"`)

	// nil => attribute omitted, so a conforming client applies the default (true).
	out, err = buildAltReplaceMPD(nil).WriteToString("  ", false)
	require.NoError(t, err)
	require.NotContains(t, out, "clip=")
}

// TestEarliestResolutionTimeOffsetFixedNotation verifies that the (large) xs:double
// earliestResolutionTimeOffset renders in plain decimal, not scientific notation
// (60s * 90000 = 5400000 previously rendered as "5.4e+06").
func TestEarliestResolutionTimeOffsetFixedNotation(t *testing.T) {
	out, err := buildAltReplaceMPD(m.Ptr(true)).WriteToString("  ", false)
	require.NoError(t, err)
	require.Contains(t, out, `earliestResolutionTimeOffset="5400000"`)
	require.NotContains(t, strings.ToLower(out), "e+0")
}
