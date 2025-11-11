package mpd_test

import (
	"bytes"
	"io/fs"
	"os"
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
	mpdCopy := m.Clone(&mpd)
	cmp.Equal(&mpd, mpdCopy)
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
