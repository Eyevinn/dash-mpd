package mpd_test

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/Eyevinn/dash-mpd/xml"
	"github.com/google/go-cmp/cmp"

	m "github.com/Eyevinn/dash-mpd/mpd"

	"aqwari.net/xml/xmltree"
	"github.com/stretchr/testify/require"
)

const (
	printProcessMPDFile = "" // Used to see the output of the process even if there is no essential diff
)

func TestDecodeEncodeMPDs(t *testing.T) {
	testDirs := []string{
		"testdata/go-dash-fixtures",
		"testdata/schema-mpds",
		"testdata/livesim",
		"testdata/hbbtv",
		"testdata/other",
		"testdata/scte35",
		"testdata/scte214",
	}
	for _, testDir := range testDirs {
		fsys := os.DirFS(testDir)
		mpdFiles, err := fs.Glob(fsys, "*.mpd")
		require.NoError(t, err)
		for _, fName := range mpdFiles {
			td, err := fs.ReadFile(fsys, fName)
			require.NoError(t, err)
			mpd, err := m.MPDFromBytes(td)
			if fName == "invalid.mpd" {
				require.Error(t, err, "")
				continue
			}
			require.NoError(t, err, fName)
			mpd.SetParents()
			out, err := xml.MarshalIndent(mpd, "", "  ")
			require.NoError(t, err)
			inTree, err := xmltree.Parse(td)
			require.NoError(t, err)
			outTree, err := xmltree.Parse(out)
			if fName == printProcessMPDFile {
				err := os.WriteFile(fName, out, 0644)
				require.NoError(t, err)
			}
			require.NoError(t, err)
			if !xmltree.Equal(inTree, outTree) {
				inBuf := bytes.Buffer{}
				err = xmltree.Encode(&inBuf, inTree)
				require.NoError(t, err)
				outBuf := bytes.Buffer{}
				err = xmltree.Encode(&outBuf, outTree)
				require.NoError(t, err)
				d := cmp.Diff(inBuf.String(), outBuf.String())
				// Note. There is no canonicalization and there may
				// be comments in the input, so the diff is not minimal.
				t.Errorf("non-minimal diff for mpd %s:\n%s. Writing file %s\n", fName, d[:400], fName)
				err = os.WriteFile(fName, out, 0644)
				require.NoError(t, err)
			}
		}
	}
}

func TestExtUrlQueryInfo(t *testing.T) {

	cases := []struct {
		name  string
		input string
	}{
		{
			name: "example_G31.1 ExtUrlQueryInfo",
			input: `<EssentialProperty schemeIdUri="urn:mpeg:dash:urlparam:2016"
	xmlns:up="urn:mpeg:dash:schema:urlparam:2016">
	<up:ExtUrlQueryInfo includeInRequests="altmpd"
		queryTemplate="prta=$urn:mpeg:dash:state:execution-delta#42$" />
</EssentialProperty>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ep := m.DescriptorType{}
			err := xml.Unmarshal([]byte(tc.input), &ep)
			require.NoError(t, err)

		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	data, err := os.ReadFile("testdata/schema-mpds/example_G15.mpd")
	require.NoError(b, err)
	for i := 0; i < b.N; i++ {
		mpd := m.MPD{}
		_ = xml.Unmarshal(data, &mpd)
	}
}

func BenchmarkMarshal(b *testing.B) {
	data, err := os.ReadFile("testdata/schema-mpds/example_G15.mpd")
	require.NoError(b, err)
	mpd := m.MPD{}
	err = xml.Unmarshal(data, &mpd)
	require.NoError(b, err)
	for i := 0; i < b.N; i++ {
		_, _ = xml.MarshalIndent(mpd, "", "  ")
	}
}

func BenchmarkClone(b *testing.B) {
	data, err := os.ReadFile("testdata/schema-mpds/example_G15.mpd")
	require.NoError(b, err)
	mpd := m.MPD{}
	err = xml.Unmarshal(data, &mpd)
	require.NoError(b, err)
	for i := 0; i < b.N; i++ {
		_ = m.Clone(&mpd)
	}
}

func TestNewFunction(t *testing.T) {

	_, err := xml.Marshal(m.NewMPD(m.DYNAMIC_TYPE))
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewPeriod())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewAdaptationSet())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewRepresentation())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewSubRepresentation())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewSegmentTemplate())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewSegmentList())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewSegmentTimeline())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewInitializationSet())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewPreselection())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewContentProtection())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewProducerReferenceTime())
	require.NoError(t, err)

	_, err = xml.Marshal(m.NewUIntVWithID())
	require.NoError(t, err)
}

func TestSegmentTemplateTimescale(t *testing.T) {
	testCases := []struct {
		timescale uint32
	}{
		{timescale: 1},
		{timescale: 1000},
	}

	for _, tc := range testCases {
		st := m.NewSegmentTemplate()
		st.SetTimescale(tc.timescale)
		gotTimescale := st.GetTimescale()
		require.Equal(t, tc.timescale, gotTimescale)
	}
}

func BenchmarkUnmarshalCorpus(b *testing.B) {
	files, err := filepath.Glob("testdata/schema-mpds/*.mpd")
	require.NoError(b, err)
	for _, f := range files {
		data, err := os.ReadFile(f)
		require.NoError(b, err)
		b.Run(filepath.Base(f), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				mpd := m.MPD{}
				_ = xml.Unmarshal(data, &mpd)
			}
		})
	}
}

func BenchmarkMarshalCorpus(b *testing.B) {
	files, err := filepath.Glob("testdata/schema-mpds/*.mpd")
	require.NoError(b, err)
	for _, f := range files {
		data, err := os.ReadFile(f)
		require.NoError(b, err)
		mpd := m.MPD{}
		err = xml.Unmarshal(data, &mpd)
		require.NoError(b, err)
		b.Run(filepath.Base(f), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, _ = xml.MarshalIndent(mpd, "", "  ")
			}
		})
	}
}

func BenchmarkWriteToDiscard(b *testing.B) {
	mpd, err := m.ReadFromFile("testdata/schema-mpds/example_G15.mpd")
	require.NoError(b, err)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = mpd.Write(io.Discard, "  ", true)
	}
}

func TestContentProtectionCerturl(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected m.ContentProtectionType
	}{
		{
			name: "ContentProtection with single Certurl",
			input: `<ContentProtection xmlns:dashif="https://dashif.org/CPS">
				<dashif:Certurl certType="application/pkcs7-mime">https://example.com/cert.pem</dashif:Certurl>
			</ContentProtection>`,
			expected: m.ContentProtectionType{
				Certurls: []m.CerturlType{
					{
						CertType: "application/pkcs7-mime",
						Value:    "https://example.com/cert.pem",
					},
				},
			},
		},
		{
			name: "ContentProtection with multiple Certurls",
			input: `<ContentProtection xmlns:dashif="https://dashif.org/CPS">
				<dashif:Certurl certType="application/pkcs7-mime">https://example.com/cert1.pem</dashif:Certurl>
				<dashif:Certurl certType="application/x-x509-ca-cert">https://example.com/cert2.pem</dashif:Certurl>
			</ContentProtection>`,
			expected: m.ContentProtectionType{
				Certurls: []m.CerturlType{
					{
						CertType: "application/pkcs7-mime",
						Value:    "https://example.com/cert1.pem",
					},
					{
						CertType: "application/x-x509-ca-cert",
						Value:    "https://example.com/cert2.pem",
					},
				},
			},
		},
		{
			name: "ContentProtection with Certurl without certType",
			input: `<ContentProtection xmlns:dashif="https://dashif.org/CPS">
				<dashif:Certurl>https://example.com/cert.pem</dashif:Certurl>
			</ContentProtection>`,
			expected: m.ContentProtectionType{
				Certurls: []m.CerturlType{
					{
						CertType: "",
						Value:    "https://example.com/cert.pem",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var cp m.ContentProtectionType
			err := xml.Unmarshal([]byte(tc.input), &cp)
			require.NoError(t, err)
			require.Equal(t, tc.expected.Certurls, cp.Certurls)
		})
	}
}

// TestEventDuration checks that a zero Event@duration can be expressed and survives a round
// trip. ISO/IEC 23009-1 gives an absent @duration the meaning "the duration is unknown", so a
// known, zero-length event is a different thing and needs to be written out. SCTE 214-1
// §6.7.2.1 relies on the difference: an SCTE-35 event that closes an earlier one "should have
// a duration of zero and shall not be infinite".
func TestEventDuration(t *testing.T) {
	cases := []struct {
		desc     string
		duration *uint64
		wantAttr string
	}{
		{desc: "unknown duration is left out", duration: nil, wantAttr: ""},
		{desc: "zero duration is written", duration: m.Ptr(uint64(0)), wantAttr: ` duration="0"`},
		{desc: "non-zero duration is written", duration: m.Ptr(uint64(90000)), wantAttr: ` duration="90000"`},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ev := &m.EventType{
				PresentationTime: 900000,
				Duration:         c.duration,
				Id:               m.Ptr(uint64(1)),
			}
			out, err := xml.Marshal(ev)
			require.NoError(t, err)
			want := `<Event presentationTime="900000"` + c.wantAttr + ` id="1"></Event>`
			require.Equal(t, want, string(out))

			var got m.EventType
			require.NoError(t, xml.Unmarshal(out, &got))
			if c.duration == nil {
				require.Nil(t, got.Duration, "an absent @duration must stay unknown")
				return
			}
			require.NotNil(t, got.Duration)
			require.Equal(t, *c.duration, *got.Duration)
		})
	}
}

func TestGroupLabel(t *testing.T) {
	const mpdStr = `<MPD xmlns="urn:mpeg:dash:schema:mpd:2011">
  <Period>
    <GroupLabel id="2">Audio</GroupLabel>
    <AdaptationSet>
      <GroupLabel id="1" lang="en">Main video</GroupLabel>
      <Label id="1" lang="sv-SE">Huvudvideo</Label>
      <Representation id="V300">
        <Label id="1" lang="en">Main video 300kbps</Label>
      </Representation>
    </AdaptationSet>
  </Period>
</MPD>`
	mpd, err := m.ReadFromString(mpdStr)
	require.NoError(t, err)
	p := mpd.Periods[0]
	require.Equal(t, []*m.LabelType{{Id: 2, Value: "Audio"}}, p.GroupLabels)
	as := p.AdaptationSets[0]
	require.Equal(t, []*m.LabelType{{Id: 1, Lang: "en", Value: "Main video"}}, as.GroupLabels)
	require.Equal(t, []*m.LabelType{{Id: 1, Lang: "sv-SE", Value: "Huvudvideo"}}, as.Labels)
	require.Equal(t, []*m.LabelType{{Id: 1, Lang: "en", Value: "Main video 300kbps"}},
		as.Representations[0].Labels)

	out, err := xml.Marshal(as)
	require.NoError(t, err)
	require.Contains(t, string(out), `<GroupLabel id="1" lang="en">Main video</GroupLabel>`)
	require.Contains(t, string(out), `<Label id="1" lang="sv-SE">Huvudvideo</Label>`)
}

// TestRepresentationBaseElementOrder checks that InitializationSet and Preselection write the
// RepresentationBase elements before Accessibility, Role, Rating and Viewpoint, as the schema requires.
func TestRepresentationBaseElementOrder(t *testing.T) {
	const initSet = `<InitializationSet id="1" contentType="video">` +
		`<EssentialProperty schemeIdUri="urn:example:feature"></EssentialProperty>` +
		`<Role schemeIdUri="urn:mpeg:dash:role:2011" value="main"></Role>` +
		`</InitializationSet>`
	const presel = `<Preselection id="1" preselectionComponents="1" lang="und">` +
		`<AudioChannelConfiguration schemeIdUri="urn:mpeg:mpegB:cicp:ChannelConfiguration" value="16">` +
		`</AudioChannelConfiguration>` +
		`<Accessibility schemeIdUri="urn:mpeg:dash:role:2011" value="description"></Accessibility>` +
		`<Role schemeIdUri="urn:mpeg:dash:role:2011" value="main"></Role>` +
		`</Preselection>`
	mpdStr := `<MPD xmlns="urn:mpeg:dash:schema:mpd:2011">` + initSet + `<Period>` + presel + `</Period></MPD>`
	mpd, err := m.ReadFromString(mpdStr)
	require.NoError(t, err)

	out, err := xml.Marshal(mpd.InitializationSet[0])
	require.NoError(t, err)
	require.Equal(t, initSet, string(out))
	out, err = xml.Marshal(mpd.Periods[0].Preselections[0])
	require.NoError(t, err)
	require.Equal(t, presel, string(out))
}

// TestEventServiceDescriptions checks that an Event keeps all its ServiceDescription elements,
// since the 6th edition schema allows more than one.
func TestEventServiceDescriptions(t *testing.T) {
	const event = `<Event presentationTime="0" id="1">` +
		`<ServiceDescription id="1"></ServiceDescription>` +
		`<ServiceDescription id="2"></ServiceDescription>` +
		`</Event>`
	var ev m.EventType
	require.NoError(t, xml.Unmarshal([]byte(event), &ev))
	require.Len(t, ev.ServiceDescriptions, 2)
	out, err := xml.Marshal(&ev)
	require.NoError(t, err)
	require.Equal(t, event, string(out))
}

// TestSegmentSequenceProperties checks that @alignment is read and written as a
// whitespace-separated list of AdaptationSet ids, and that event="false" is kept, since a
// missing @event means true.
func TestSegmentSequenceProperties(t *testing.T) {
	const ssp = `<SegmentSequenceProperties sapType="1" cadence="4" event="false" alignment="1 2">` +
		`</SegmentSequenceProperties>`
	mpdStr := `<MPD xmlns="urn:mpeg:dash:schema:mpd:2011"><Period><AdaptationSet id="1">` + ssp +
		`</AdaptationSet></Period></MPD>`
	mpd, err := m.ReadFromString(mpdStr)
	require.NoError(t, err)
	as := mpd.Periods[0].AdaptationSets[0]
	props := as.SegmentSequenceProperties[0]
	require.Equal(t, &m.StringNoWhitespaceVectorType{"1", "2"}, props.Alignment)
	require.Equal(t, m.Ptr(false), props.Event)
	out, err := xml.Marshal(as)
	require.NoError(t, err)
	require.Contains(t, string(out), ssp)

	// nil => attribute omitted, so a conforming client applies the default (true).
	props.Event = nil
	out, err = xml.Marshal(as)
	require.NoError(t, err)
	require.NotContains(t, string(out), "event=")
}
