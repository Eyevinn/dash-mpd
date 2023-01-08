package mpd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/barkimedes/go-deepcopy"
)

// Clone creates a deep copy of mpd.
func Clone(mpd *MPD) *MPD {
	return deepcopy.MustAnything(mpd).(*MPD)
}

// GetRepInit returns the representation's initialization URI with replaced identifiers.

// BaseURLs are not applied.
func GetRepInit(a *AdaptationSetType, r *RepresentationType) (string, error) {
	if a == nil || r == nil {
		return "", fmt.Errorf("need both adaptationSet and representation")
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

// BaseURLs are not applied.
func GetRepMedia(a *AdaptationSetType, r *RepresentationType) (string, error) {
	if a == nil || r == nil {
		return "", fmt.Errorf("need both adaptationSet and representation")
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
