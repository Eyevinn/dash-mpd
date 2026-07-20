package mpd_test

import (
	"os"
	"testing"

	"github.com/Eyevinn/dash-mpd/mpd"
	"github.com/Eyevinn/dash-mpd/xml"
	"github.com/stretchr/testify/require"
)

// TestSCTE214SupplementalAttributes verifies that the SCTE-214
// supplementalProfiles and supplementalCodecs attributes are read on both
// AdaptationSet and Representation, whatever prefix the input binds to the
// SCTE-214 namespace URI.
func TestSCTE214SupplementalAttributes(t *testing.T) {
	cases := []struct {
		desc         string
		file         string
		wantProfiles []mpd.StringVectorType
		wantCodecs   []mpd.StringVectorType
	}{
		{
			desc:         "Representation level, one value per attribute",
			file:         "testdata/scte214/supplemental_codecs_representation.mpd",
			wantProfiles: []mpd.StringVectorType{"db1p", "db1p"},
			wantCodecs:   []mpd.StringVectorType{"dvh1.08.01", "dvh1.08.03"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			data, err := os.ReadFile(tc.file)
			require.NoError(t, err)
			m, err := mpd.MPDFromBytes(data)
			require.NoError(t, err)

			reps := m.Periods[0].AdaptationSets[0].Representations
			require.Len(t, reps, len(tc.wantCodecs))
			for i, rep := range reps {
				require.Equal(t, tc.wantProfiles[i], rep.SupplementalProfiles)
				require.Equal(t, tc.wantCodecs[i], rep.SupplementalCodecs)
			}
		})
	}
}

// TestSCTE214SupplementalAttributesOnAdaptationSet covers the AdaptationSet
// placement and space-separated multi-value lists.
func TestSCTE214SupplementalAttributesOnAdaptationSet(t *testing.T) {
	data, err := os.ReadFile("testdata/scte214/supplemental_codecs_adaptation_set.mpd")
	require.NoError(t, err)
	m, err := mpd.MPDFromBytes(data)
	require.NoError(t, err)

	as := m.Periods[0].AdaptationSets[0]
	require.Equal(t, mpd.StringVectorType("db1p cdm4"), as.SupplementalProfiles)
	require.Equal(t, mpd.StringVectorType("dvh1.08.09"), as.SupplementalCodecs)
}

// TestSCTE214NamespacePrefixInsensitive verifies that the attributes are
// matched by namespace URI, not by prefix, so manifests binding any prefix
// to the SCTE-214 namespace are read correctly.
func TestSCTE214NamespacePrefixInsensitive(t *testing.T) {
	input := `<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" xmlns:_="urn:scte:dash:scte214-extensions"` +
		` profiles="urn:mpeg:dash:profile:isoff-on-demand:2011" type="static"` +
		` mediaPresentationDuration="PT30S" minBufferTime="PT2S">` +
		`<Period id="1"><AdaptationSet id="1" mimeType="video/mp4">` +
		`<Representation id="v0" bandwidth="4725303" codecs="hvc1.2.4.L150.90"` +
		` _:supplementalProfiles="db1p" _:supplementalCodecs="dvh1.08.05"></Representation>` +
		`</AdaptationSet></Period></MPD>`

	m, err := mpd.MPDFromBytes([]byte(input))
	require.NoError(t, err)
	rep := m.Periods[0].AdaptationSets[0].Representations[0]
	require.Equal(t, mpd.StringVectorType("db1p"), rep.SupplementalProfiles)
	require.Equal(t, mpd.StringVectorType("dvh1.08.05"), rep.SupplementalCodecs)

	// On marshal, the attributes come out with the scte214 prefix declared
	// on the carrying element.
	out, err := xml.MarshalIndent(m, "", "  ")
	require.NoError(t, err)
	s := string(out)
	require.Contains(t, s, `xmlns:scte214="urn:scte:dash:scte214-extensions"`)
	require.Contains(t, s, `scte214:supplementalProfiles="db1p"`)
	require.Contains(t, s, `scte214:supplementalCodecs="dvh1.08.05"`)
}

// TestSCTE214ContentIdentifiers covers ContentIdentifier elements inside
// AssetIdentifier and SupplementalProperty descriptors with the
// urn:scte:dash:asset-id:upid:2015 scheme.
func TestSCTE214ContentIdentifiers(t *testing.T) {
	data, err := os.ReadFile("testdata/scte214/asset_id_upid.mpd")
	require.NoError(t, err)
	m, err := mpd.MPDFromBytes(data)
	require.NoError(t, err)

	p := m.Periods[0]
	ai := p.AssetIdentifier
	require.NotNil(t, ai)
	require.Equal(t, mpd.AnyURI(mpd.SCTE214SchemeIdAssetIdUpid), ai.SchemeIdUri)
	require.Len(t, ai.ContentIdentifiers, 2)
	require.Equal(t, "ADI", ai.ContentIdentifiers[0].Type)
	require.Equal(t, "cablelabs.com/MOVE1234567890123456", ai.ContentIdentifiers[0].Value)
	require.Equal(t, "MPU", ai.ContentIdentifiers[1].Type)
	require.Equal(t, "CSP1DE12AB327FE312AF", ai.ContentIdentifiers[1].Value)

	require.Len(t, p.SupplementalProperties, 1)
	sp := p.SupplementalProperties[0]
	require.Equal(t, mpd.AnyURI(mpd.SCTE214SchemeIdAssetIdUpid), sp.SchemeIdUri)
	require.Len(t, sp.ContentIdentifiers, 1)
	require.Equal(t, "AiringID", sp.ContentIdentifiers[0].Type)
	require.Equal(t, "0xDEADBEEF", sp.ContentIdentifiers[0].Value)
}

// TestSCTE214ProgrammaticConstruction builds the SCTE-214 signalling
// programmatically, marshals it and checks the output, then parses it back.
func TestSCTE214ProgrammaticConstruction(t *testing.T) {
	m := mpd.NewMPD(mpd.STATIC_TYPE)
	m.Profiles = mpd.PROFILE_ONDEMAND
	m.MediaPresentationDuration = mpd.Seconds2DurPtr(30)
	m.MinBufferTime = mpd.Seconds2DurPtr(2)

	p := mpd.NewPeriod()
	p.Id = "1"
	m.AppendPeriod(p)

	ai := mpd.NewDescriptor(mpd.SCTE214SchemeIdAssetIdUpid, "", "")
	ai.ContentIdentifiers = []*mpd.ContentIdentifierType{
		mpd.NewContentIdentifier("EIDR", "10.5240/EA73-79D7-1B2B-B378-3A73-M"),
	}
	p.AssetIdentifier = ai

	as := mpd.NewAdaptationSet()
	as.MimeType = "video/mp4"
	p.AppendAdaptationSet(as)

	rep := mpd.NewRepresentation()
	rep.Id = "video"
	rep.Bandwidth = 5139658
	rep.Codecs = "hvc1.2.4.L120.b0"
	rep.SupplementalProfiles = "db1p cdm4"
	rep.SupplementalCodecs = "dvh1.08.03"
	as.AppendRepresentation(rep)

	out, err := xml.MarshalIndent(m, "", "  ")
	require.NoError(t, err)
	s := string(out)
	require.Contains(t, s, `xmlns:scte214="urn:scte:dash:scte214-extensions"`)
	require.Contains(t, s, `scte214:supplementalProfiles="db1p cdm4"`)
	require.Contains(t, s, `scte214:supplementalCodecs="dvh1.08.03"`)
	require.Contains(t, s, `<scte214:ContentIdentifier`)
	require.Contains(t, s, `type="EIDR" value="10.5240/EA73-79D7-1B2B-B378-3A73-M"`)

	rt, err := mpd.MPDFromBytes(out)
	require.NoError(t, err)
	rtRep := rt.Periods[0].AdaptationSets[0].Representations[0]
	require.Equal(t, mpd.StringVectorType("db1p cdm4"), rtRep.SupplementalProfiles)
	require.Equal(t, mpd.StringVectorType("dvh1.08.03"), rtRep.SupplementalCodecs)
	rtAi := rt.Periods[0].AssetIdentifier
	require.NotNil(t, rtAi)
	require.Len(t, rtAi.ContentIdentifiers, 1)
	require.Equal(t, "EIDR", rtAi.ContentIdentifiers[0].Type)
}
