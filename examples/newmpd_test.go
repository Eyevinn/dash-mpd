package mpd_test

import (
	"fmt"

	"github.com/Eyevinn/dash-mpd/xml"

	"github.com/Eyevinn/dash-mpd/mpd"
)

func ExampleNewMPD() {
	m := mpd.NewMPD(mpd.STATIC_TYPE)
	m.Profiles = "urn:mpeg:dash:profile:isoff-live:2011,http://dashif.org/guidelines/dash-if-simple"
	p := mpd.NewPeriod()
	m.AppendPeriod(p)
	p.Id = "p0"
	as := mpd.NewAdaptationSet()
	p.AppendAdaptationSet(as)
	as.ContentType = "audio"
	as.Lang = "en"
	st := mpd.NewSegmentTemplate()
	as.SegmentTemplate = st
	st.StartNumber = mpd.Ptr(uint32(1))
	st.Initialization = "$RepresentationID$/init.mp4"
	st.Duration = mpd.Ptr(uint32(2))
	st.Media = "$RepresentationID$/$Number$.m4s"
	rep := mpd.NewRepresentation()
	as.AppendRepresentation(rep)
	rep.Id = "A48"
	rep.Codecs = "mp4a.40.2"
	rep.Bandwidth = 96000
	rep.AudioSamplingRate = mpd.Ptr(mpd.UIntVectorType("48000"))
	out, _ := xml.MarshalIndent(m, " ", "")

	fmt.Println(string(out))
	//nolint:lll
	// Output: <MPD xmlns="urn:mpeg:dash:schema:mpd:2011" profiles="urn:mpeg:dash:profile:isoff-live:2011,http://dashif.org/guidelines/dash-if-simple" type="static">
	//  <Period id="p0">
	//  <AdaptationSet lang="en" contentType="audio">
	//  <SegmentTemplate media="$RepresentationID$/$Number$.m4s" initialization="$RepresentationID$/init.mp4" duration="2" startNumber="1"></SegmentTemplate>
	//  <Representation id="A48" bandwidth="96000" audioSamplingRate="48000" codecs="mp4a.40.2"></Representation>
	//  </AdaptationSet>
	//  </Period>
	//  </MPD>
}
