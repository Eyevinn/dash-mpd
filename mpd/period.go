package mpd

import "fmt"

// AbsoluteStart returns the absolute start time of a Period in seconds.
// It is meant for dynamic MPDs where each Period has a start time,
// but it can also be used for the start-over case where a dynamice MPD
// has just turned into a static MPD and still have availabilityStartTime,
// with the same value as in the dynamice case.
func (p *Period) AbsoluteStart(m *MPD) (float64, error) {
	if m.AvailabilityStartTime == "" {
		return 0, ErrASTRequired
	}
	ast, err := m.AvailabilityStartTime.ConvertToSeconds()
	if err != nil {
		return 0, err
	}
	pStart, err := p.GetStart()
	if err != nil {
		return 0, err
	}
	return ast + pStart.Seconds(), nil
}

func (p *Period) GetIndex() (int, error) {
	m := p.Parent()
	if m == nil {
		return -1, ErrParentNotSet
	}
	for i, period := range m.Periods {
		if period == p {
			return i, nil
		}
	}
	return -1, ErrPeriodNotFound
}

// PeriodType is the type of period
type PeriodType int

const (
	PTBad PeriodType = iota
	PTRegular
	PTEarlyAvailable
	PTEarlyTerminated
	PTRegularOrEarlyTerminated
)

func (p PeriodType) String() string {
	switch p {
	case PTRegular:
		return "regular"
	case PTEarlyAvailable:
		return "early-available"
	case PTEarlyTerminated:
		return "early-terminated"
	}
	return "unknown"
}

// GetType returns the type of period.
func (p *Period) GetType() (PeriodType, error) {
	m := p.Parent()
	if m == nil {
		return PTBad, ErrParentNotSet
	}
	nrPeriods := len(m.Periods)
	idx, err := p.GetIndex()
	if err != nil {
		return PTBad, err
	}
	isEarlyTerminated := func(p *Period, m *MPD, idx, nrPeriods int) bool {
		return p.Duration != nil && (m.MinimumUpdatePeriod != nil ||
			(idx < nrPeriods-1 && m.Periods[idx+1].Start != nil))
	}
	if p.Start != nil {
		if isEarlyTerminated(p, m, idx, nrPeriods) {
			return PTEarlyTerminated, nil
		}
		return PTRegular, nil
	}
	// No start attribute
	if idx > 0 && m.Periods[idx-1].Duration != nil {
		// Regular or Early Terminated
		if isEarlyTerminated(p, m, idx, nrPeriods) {
			return PTEarlyTerminated, nil
		}
		return PTRegularOrEarlyTerminated, nil
	}
	if m.GetType() == "dynamic" && (idx == 0 || m.Periods[idx-1].Duration == nil) {
		return PTEarlyAvailable, nil
	}
	return PTRegular, nil
}

// GetStart returns the period's start time in seconds (relative to availabilityStartTime if present).
// The algorithm is specified in ISO/IEC 23009-1 Section 5.3.2.1.
func (p *Period) GetStart() (Duration, error) {
	m := p.Parent()
	if m == nil {
		return 0, ErrParentNotSet
	}
	if p.Start != nil {
		return *p.Start, nil
	}
	idx, err := p.GetIndex()
	if err != nil {
		return 0, err
	}
	if idx > 0 {
		prev := m.Periods[idx-1]
		if prev.Duration != nil {
			prevStart, err := prev.GetStart()
			if err != nil {
				return 0, fmt.Errorf("cannot get start of previous period: %w", err)
			}
			return prevStart + *prev.Duration, nil
		}
		return 0, ErrUnknownPeriodStart
	}
	// first period of static, but no start => start = 0
	if m.GetType() == StaticMPDType {
		return 0, nil
	}
	return 0, ErrUnknownPeriodStart
}

// GetDuration returns the period's duration if specified or the MPD is single-period static.
func (p *Period) GetDuration() (Duration, error) {
	if p.Duration != nil {
		return *p.Duration, nil
	}
	m := p.Parent()
	if m == nil {
		return 0, ErrParentNotSet
	}
	nrPeriods := len(m.Periods)
	idx, err := p.GetIndex()
	if err != nil {
		return 0, err
	}
	switch m.GetType() {
	case StaticMPDType:
		if m.MediaPresentationDuration == nil {
			return 0, ErrNoMediaPresentationDuration
		}
		end := *m.MediaPresentationDuration
		if idx < nrPeriods-1 {
			if m.Periods[idx+1].Start == nil {
				return 0, ErrUnknownPeriodDur
			}
			end = *m.Periods[idx+1].Start
		}
		if p.Start == nil {
			return 0, ErrUnknownPeriodDur
		}
		return end - *p.Start, nil
	default: // dynamic
		if p.Start == nil {
			// Early available period
			return 0, ErrUnknownPeriodDur
		}
		if idx == nrPeriods-1 {
			return 0, ErrUnknownPeriodDur
		}
		if m.Periods[idx+1].Start == nil {
			return 0, ErrUnknownPeriodDur
		}
		return *m.Periods[idx+1].Start - *p.Start, nil
	}
}
