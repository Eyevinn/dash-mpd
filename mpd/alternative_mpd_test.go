package mpd_test

import (
	"testing"

	m "github.com/Eyevinn/dash-mpd/mpd"
	"github.com/Eyevinn/dash-mpd/xml"
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
					EarliestResolutionTimeOffset: m.Ptr(uint64(5400000)),
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

// TestAlternativeMPDEventAttributes verifies that earliestResolutionTimeOffset, which has no
// default, is left out when nil but written when 0, that a large value renders in plain decimal
// (5400000 once rendered as "5.4e+06"), and that serviceDescriptionId round trips.
func TestAlternativeMPDEventAttributes(t *testing.T) {
	cases := []struct {
		desc     string
		offset   *uint64
		wantAttr string
	}{
		{desc: "absent offset is left out", offset: nil, wantAttr: ""},
		{desc: "zero offset is written", offset: m.Ptr(uint64(0)), wantAttr: ` earliestResolutionTimeOffset="0"`},
		{desc: "large offset is plain decimal", offset: m.Ptr(uint64(5400000)), wantAttr: ` earliestResolutionTimeOffset="5400000"`},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ev := &m.EventType{
				InsertPresentation: &m.AlternativeMPDEventType{
					Uri:                          "http://example.com/ad.mpd",
					EarliestResolutionTimeOffset: c.offset,
					ServiceDescriptionId:         m.Ptr(uint32(2)),
				},
			}
			out, err := xml.Marshal(ev)
			require.NoError(t, err)
			want := `<Event presentationTime="0"><InsertPresentation uri="http://example.com/ad.mpd"` +
				c.wantAttr + ` serviceDescriptionId="2"></InsertPresentation></Event>`
			require.Equal(t, want, string(out))

			var got m.EventType
			require.NoError(t, xml.Unmarshal(out, &got))
			require.Equal(t, ev.InsertPresentation, got.InsertPresentation)
		})
	}
}
