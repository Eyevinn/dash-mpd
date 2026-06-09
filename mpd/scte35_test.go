package mpd_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Eyevinn/dash-mpd/mpd"
	"github.com/Eyevinn/dash-mpd/xml"
	"github.com/stretchr/testify/require"
)

// TestSCTE35NamespacePassThrough confirms that the SCTE-35 namespace URI
// observed at unmarshal time is replayed verbatim at marshal time, for both
// the canonical "http://www.scte.org/schemas/35" URI and the legacy 2016
// variant used in much of the install base.
func TestSCTE35NamespacePassThrough(t *testing.T) {
	cases := []struct {
		file     string
		wantNS   string
		wantLine string
	}{
		{
			file:     "testdata/scte35/time_signal_segmentation.mpd",
			wantNS:   mpd.SCTE35Namespace,
			wantLine: `<Signal xmlns="http://www.scte.org/schemas/35">`,
		},
		{
			file:     "testdata/scte35/signal_binary_xmlbin.mpd",
			wantNS:   mpd.SCTE35Namespace2016,
			wantLine: `<Signal xmlns="http://www.scte.org/schemas/35/2016">`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			data, err := os.ReadFile(tc.file)
			require.NoError(t, err)
			m, err := mpd.MPDFromBytes(data)
			require.NoError(t, err)

			ev := m.Periods[0].EventStreams[0].Events[0]
			require.NotNil(t, ev.Signal, "expected Signal to be populated")
			require.Equal(t, tc.wantNS, ev.Signal.XMLName.Space, "namespace not captured from input")

			out, err := xml.MarshalIndent(m, "", "  ")
			require.NoError(t, err)
			require.Contains(t, string(out), tc.wantLine,
				"marshal output does not preserve input namespace")
		})
	}
}

// TestSCTE35SpliceInfoSectionDirectFormPassThrough is the equivalent of
// TestSCTE35NamespacePassThrough for the "no Signal wrapper" form
// (SpliceInfoSection placed directly under Event, AWS MediaTailor style).
func TestSCTE35SpliceInfoSectionDirectFormPassThrough(t *testing.T) {
	data, err := os.ReadFile("testdata/scte35/direct_splice_info_section.mpd")
	require.NoError(t, err)
	m, err := mpd.MPDFromBytes(data)
	require.NoError(t, err)

	ev := m.Periods[0].EventStreams[0].Events[0]
	require.Nil(t, ev.Signal, "Signal must be nil for the direct form")
	require.NotNil(t, ev.SpliceInfoSection, "SpliceInfoSection must be populated for the direct form")
	require.Equal(t, mpd.SCTE35Namespace, ev.SpliceInfoSection.XMLName.Space)

	out, err := xml.MarshalIndent(m, "", "  ")
	require.NoError(t, err)
	s := string(out)
	require.Contains(t, s, `<SpliceInfoSection`,
		"SpliceInfoSection element must be emitted")
	require.Contains(t, s, mpd.SCTE35Namespace,
		"direct SpliceInfoSection should round-trip with its namespace preserved")
}

// TestSCTE35SpliceInfo verifies the (*EventType).SpliceInfo helper returns
// the section from whichever form (wrapped or direct) is populated.
func TestSCTE35SpliceInfo(t *testing.T) {
	t.Run("wrapped in Signal", func(t *testing.T) {
		data, err := os.ReadFile("testdata/scte35/time_signal_segmentation.mpd")
		require.NoError(t, err)
		m, err := mpd.MPDFromBytes(data)
		require.NoError(t, err)

		sis := m.Periods[0].EventStreams[0].Events[0].SpliceInfo()
		require.NotNil(t, sis)
		require.NotNil(t, sis.TimeSignal)
	})
	t.Run("direct under Event", func(t *testing.T) {
		data, err := os.ReadFile("testdata/scte35/direct_splice_info_section.mpd")
		require.NoError(t, err)
		m, err := mpd.MPDFromBytes(data)
		require.NoError(t, err)

		sis := m.Periods[0].EventStreams[0].Events[0].SpliceInfo()
		require.NotNil(t, sis)
		require.NotNil(t, sis.TimeSignal)
	})
	t.Run("neither set", func(t *testing.T) {
		ev := &mpd.EventType{}
		require.Nil(t, ev.SpliceInfo())
	})
}

// TestSCTE35ConstructorsUseCanonicalNamespace asserts that the New* helpers
// in mpd/scte35.go produce values whose XMLName carries the canonical
// SCTE-35 URI, so that programmatic construction marshals to a proper
// SCTE-35 element by default.
func TestSCTE35ConstructorsUseCanonicalNamespace(t *testing.T) {
	cases := []struct {
		name  string
		got   xml.Name
		local string
	}{
		{"Signal", mpd.NewSignal().XMLName, "Signal"},
		{"SpliceInfoSection", mpd.NewSpliceInfoSection().XMLName, "SpliceInfoSection"},
		{"Binary", mpd.NewBinary("").XMLName, "Binary"},
		{"SpliceNull", mpd.NewSpliceNull().XMLName, "SpliceNull"},
		{"SpliceSchedule", mpd.NewSpliceSchedule().XMLName, "SpliceSchedule"},
		{"SpliceInsert", mpd.NewSpliceInsert().XMLName, "SpliceInsert"},
		{"TimeSignal", mpd.NewTimeSignal().XMLName, "TimeSignal"},
		{"BandwidthReservation", mpd.NewBandwidthReservation().XMLName, "BandwidthReservation"},
		{"PrivateCommand", mpd.NewPrivateCommand(0).XMLName, "PrivateCommand"},
		{"SpliceTime", mpd.NewSpliceTime().XMLName, "SpliceTime"},
		{"BreakDuration", mpd.NewBreakDuration(false, 0).XMLName, "BreakDuration"},
		{"SegmentationDescriptor", mpd.NewSegmentationDescriptor().XMLName, "SegmentationDescriptor"},
		{"SegmentationUpid", mpd.NewSegmentationUpid("").XMLName, "SegmentationUpid"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, mpd.SCTE35Namespace, tc.got.Space)
			require.Equal(t, tc.local, tc.got.Local)
		})
	}
}

// TestEventStreamIsSCTE35 covers the schemeIdUri detection helper.
func TestEventStreamIsSCTE35(t *testing.T) {
	cases := []struct {
		scheme string
		want   bool
	}{
		{mpd.SCTE35SchemeIdXML, true},
		{mpd.SCTE35SchemeIdXMLBin, true},
		{mpd.SCTE35SchemeIdBin, true},
		{"urn:mpeg:dash:event:2012", false},
		{"", false},
	}
	for _, tc := range cases {
		t.Run(tc.scheme, func(t *testing.T) {
			es := &mpd.EventStreamType{SchemeIdUri: mpd.AnyURI(tc.scheme)}
			require.Equal(t, tc.want, es.IsSCTE35())
		})
	}
}

// TestSCTE35ProgrammaticRoundTrip builds a Signal/SpliceInfoSection/TimeSignal
// programmatically, marshals it, unmarshals back and checks key fields. Also
// asserts that the canonical SCTE-35 namespace is used in the output, even
// when wrapped inside a DASH MPD whose default namespace is the DASH one.
func TestSCTE35ProgrammaticRoundTrip(t *testing.T) {
	m := mpd.NewMPD(mpd.DYNAMIC_TYPE)
	m.Profiles = mpd.PROFILE_LIVE
	m.AvailabilityStartTime = mpd.DateTime("2024-01-01T00:00:00Z")
	m.MinBufferTime = mpd.Seconds2DurPtr(4)
	m.MinimumUpdatePeriod = mpd.Seconds2DurPtr(2)

	p := mpd.NewPeriod()
	p.Id = "1"
	p.Start = mpd.Seconds2DurPtr(0)
	m.AppendPeriod(p)

	ts := uint32(90000)
	es := &mpd.EventStreamType{
		SchemeIdUri: mpd.AnyURI(mpd.SCTE35SchemeIdXML),
		Timescale:   &ts,
	}
	p.EventStreams = append(p.EventStreams, es)

	sig := mpd.NewSignal()
	sis := mpd.NewSpliceInfoSection()
	tier := uint16(4095)
	sis.Tier = &tier

	st := mpd.NewSpliceTime()
	pts := uint64(123456789)
	st.PtsTime = &pts
	sis.TimeSignal = &mpd.TimeSignalType{
		XMLName:    xml.Name{Space: mpd.SCTE35Namespace, Local: "TimeSignal"},
		SpliceTime: st,
	}

	upid := mpd.NewSegmentationUpid("PROVIDERX")
	upidType := uint8(9)
	upid.SegmentationUpidType = &upidType

	segDesc := mpd.NewSegmentationDescriptor()
	eventID := uint32(42)
	dur := uint64(2700000)
	segType := mpd.SegTypeProviderPlacementOpportunityStart
	segNum := uint8(1)
	expected := uint8(1)
	segDesc.SegmentationEventId = &eventID
	segDesc.SegmentationDuration = &dur
	segDesc.SegmentationTypeId = &segType
	segDesc.SegmentNum = &segNum
	segDesc.SegmentsExpected = &expected
	segDesc.SegmentationUpids = []*mpd.SegmentationUpidType{upid}
	sis.SegmentationDescriptors = []*mpd.SegmentationDescriptorType{segDesc}

	sig.SpliceInfoSection = sis

	id := uint64(7)
	ev := &mpd.EventType{
		PresentationTime: 900000,
		Duration:         2700000,
		Id:               &id,
		Signal:           sig,
	}
	es.Events = append(es.Events, ev)

	out, err := xml.MarshalIndent(m, "", "  ")
	require.NoError(t, err)
	s := string(out)
	require.Contains(t, s, `<Signal xmlns="http://www.scte.org/schemas/35">`)
	require.Contains(t, s, `<TimeSignal>`)
	require.Contains(t, s, `<SpliceTime ptsTime="123456789">`)
	require.Contains(t, s, `<SegmentationDescriptor`)
	require.Contains(t, s, `segmentationTypeId="52"`)
	require.Contains(t, s, `<SegmentationUpid segmentationUpidType="9">PROVIDERX</SegmentationUpid>`)
	require.NotContains(t, s, `xmlns="http://www.scte.org/schemas/35"`+"\n",
		"namespace should be declared once, on the Signal element")

	// Round-trip back: parse the marshal output and assert structure.
	rt, err := mpd.MPDFromBytes(out)
	require.NoError(t, err)
	rtEv := rt.Periods[0].EventStreams[0].Events[0]
	require.NotNil(t, rtEv.Signal)
	require.NotNil(t, rtEv.Signal.SpliceInfoSection)
	require.NotNil(t, rtEv.Signal.SpliceInfoSection.TimeSignal)
	require.NotNil(t, rtEv.Signal.SpliceInfoSection.TimeSignal.SpliceTime)
	require.NotNil(t, rtEv.Signal.SpliceInfoSection.TimeSignal.SpliceTime.PtsTime)
	require.Equal(t, uint64(123456789), *rtEv.Signal.SpliceInfoSection.TimeSignal.SpliceTime.PtsTime)
	require.Len(t, rtEv.Signal.SpliceInfoSection.SegmentationDescriptors, 1)
	rtSeg := rtEv.Signal.SpliceInfoSection.SegmentationDescriptors[0]
	require.Equal(t, uint8(52), *rtSeg.SegmentationTypeId)
	require.Len(t, rtSeg.SegmentationUpids, 1)
	require.Equal(t, "PROVIDERX", rtSeg.SegmentationUpids[0].Value)
}

// TestSCTE35BareConstructionFallsBackToTypeName documents the behaviour when
// users build a SignalType with a struct literal and forget to set XMLName:
// since the xml fork falls through to reflect.Type.Name() when no XMLName
// and no parent field tag give a local name, the element is emitted with the
// raw Go type name ("SignalType"). Always prefer NewSignal() (and the other
// New* helpers) to avoid this surprise.
func TestSCTE35BareConstructionFallsBackToTypeName(t *testing.T) {
	bare := &mpd.SignalType{Binary: mpd.NewBinary("AA==")}
	out, err := xml.Marshal(bare)
	require.NoError(t, err)
	s := string(out)
	require.True(t, strings.HasPrefix(s, "<SignalType>"),
		"bare Signal falls back to the Go type name; got %q", s)

	// With the constructor, the element is named "Signal" and carries the
	// canonical SCTE-35 namespace.
	good := mpd.NewSignal()
	good.Binary = mpd.NewBinary("AA==")
	outGood, err := xml.Marshal(good)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(string(outGood),
		`<Signal xmlns="http://www.scte.org/schemas/35">`),
		"NewSignal() should produce a proper Signal element; got %q", outGood)
}

// TestSCTE35SpliceInsertFields covers SpliceInsert with the Program/SpliceTime
// path including BreakDuration; uses the splice_insert.mpd fixture.
func TestSCTE35SpliceInsertFields(t *testing.T) {
	data, err := os.ReadFile("testdata/scte35/splice_insert.mpd")
	require.NoError(t, err)
	m, err := mpd.MPDFromBytes(data)
	require.NoError(t, err)

	sis := m.Periods[0].EventStreams[0].Events[0].SpliceInfo()
	require.NotNil(t, sis)
	si := sis.SpliceInsert
	require.NotNil(t, si)
	require.Equal(t, uint32(4157), *si.SpliceEventId)
	require.True(t, *si.OutOfNetworkIndicator)
	require.False(t, *si.SpliceImmediateFlag)
	require.NotNil(t, si.Program)
	require.NotNil(t, si.Program.SpliceTime)
	require.Equal(t, uint64(2346545680), *si.Program.SpliceTime.PtsTime)
	require.NotNil(t, si.BreakDuration)
	require.True(t, si.BreakDuration.AutoReturn)
	require.Equal(t, uint64(2699769), si.BreakDuration.Duration)
}

// TestSCTE35BinaryRoundTrip covers the xml+bin form (Signal + Binary) end to
// end on the signal_binary_xmlbin.mpd fixture, including the legacy 2016
// namespace.
func TestSCTE35BinaryRoundTrip(t *testing.T) {
	data, err := os.ReadFile("testdata/scte35/signal_binary_xmlbin.mpd")
	require.NoError(t, err)
	m, err := mpd.MPDFromBytes(data)
	require.NoError(t, err)

	ev := m.Periods[0].EventStreams[0].Events[0]
	require.NotNil(t, ev.Signal)
	require.NotNil(t, ev.Signal.Binary)
	require.Nil(t, ev.Signal.SpliceInfoSection)
	require.Contains(t, ev.Signal.Binary.Value, "/DBI")

	out, err := xml.MarshalIndent(m, "", "  ")
	require.NoError(t, err)
	s := string(out)
	require.Contains(t, s, mpd.SCTE35Namespace2016)
	require.NotContains(t, s, `xmlns="http://www.scte.org/schemas/35">`+"\n",
		"canonical namespace must not appear when input is in the legacy URI")
}

// TestSCTE35DeliveryRestrictions verifies the DeliveryRestrictions sub-element
// of SegmentationDescriptor is parsed with all four required boolean/uint
// attributes.
func TestSCTE35DeliveryRestrictions(t *testing.T) {
	data, err := os.ReadFile("testdata/scte35/delivery_restrictions.mpd")
	require.NoError(t, err)
	m, err := mpd.MPDFromBytes(data)
	require.NoError(t, err)

	sis := m.Periods[0].EventStreams[0].Events[0].SpliceInfo()
	require.NotNil(t, sis)
	require.Len(t, sis.SegmentationDescriptors, 1)
	dr := sis.SegmentationDescriptors[0].DeliveryRestrictions
	require.NotNil(t, dr)
	require.True(t, dr.WebDeliveryAllowedFlag)
	require.True(t, dr.NoRegionalBlackoutFlag)
	require.True(t, dr.ArchiveAllowedFlag)
	require.Equal(t, uint8(3), dr.DeviceRestrictions)
}
