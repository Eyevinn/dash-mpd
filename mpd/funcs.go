package mpd

import (
	"strconv"
	"strings"

	"github.com/barkimedes/go-deepcopy"
)

// ParseDurationError for parsing xs:Duration string
type MPDError struct {
	msg string
}

func (m MPDError) Error() string {
	return m.msg
}

func newMPDError(msg string) MPDError {
	return MPDError{msg: msg}
}

// Clone creates a deep copy of mpd.
func Clone(mpd *MPD) *MPD {
	return deepcopy.MustAnything(mpd).(*MPD)
}

// GetType returns static or dynamic.
func (m *MPD) GetType() string {
	if m.Type == nil {
		return "static"
	}
	return *m.Type
}

// GetRepInit returns the representation's initialization URI with replaced identifiers.
//
// BaseURLs are not applied.
func GetRepInit(a *AdaptationSetType, r *RepresentationType) (string, error) {
	if a == nil || r == nil {
		return "", newMPDError("need both adaptationSet and representation")
	}
	var initialization string
	if a.SegmentTemplate != nil {
		initialization = a.SegmentTemplate.Initialization
	}
	if r.SegmentTemplate != nil {
		initialization = r.SegmentTemplate.Initialization
	}
	initialization = strings.ReplaceAll(initialization, "$RepresentationID$", r.Id)
	initialization = strings.ReplaceAll(initialization, "$Bandwidth$", strconv.Itoa(int(r.Bandwidth)))
	return initialization, nil
}

// GetRepMedia returns the representaion's media path with replaced ID and bandwidth identifiers.
//
// BaseURLs are not applied.
func GetRepMedia(a *AdaptationSetType, r *RepresentationType) (string, error) {
	if a == nil || r == nil {
		return "", newMPDError("need both adaptationSet and representation")
	}
	var media string
	if a.SegmentTemplate != nil {
		media = a.SegmentTemplate.Media
	}
	if r.SegmentTemplate != nil {
		media = r.SegmentTemplate.Media
	}
	media = strings.ReplaceAll(media, "$RepresentationID$", r.Id)
	media = strings.ReplaceAll(media, "$Bandwidth$", strconv.Itoa(int(r.Bandwidth)))

	return media, nil
}

// GetPeriodDuration returns the period's duration if specified or the MPD is single-period static.
func GetPeriodDuration(mpd *MPD, per *PeriodType) (Duration, error) {
	if per.Duration != nil {
		return *per.Duration, nil
	}
	if mpd.GetType() == "dynamic" || len(mpd.Periods) != 1 {
		return 0, newMPDError("cannot determine duration for dynamic or multi-period MPD")
	}
	return *mpd.MediaPresentationDuration, nil
}
