package mpd

import "errors"

var (
	ErrASTRequired                 = errors.New("availabilityStartTime is required for dynamic MPDs")
	ErrPeriodNotFound              = errors.New("period not found in MPD")
	ErrParentNotSet                = errors.New("parent not set")
	ErrNoMediaPresentationDuration = errors.New("no MediaPresentationDuration in static MPD")
	ErrUnknownPeriodDur            = errors.New("period duration cannot be derived")
	ErrNoStartInDynamicPeriod      = errors.New("start is required for dynamic periods")
	ErrUnknownPeriodStart          = errors.New("period start cannot be derived")
)

type MPDError struct {
	msg string
}

func (m MPDError) Error() string {
	return m.msg
}
