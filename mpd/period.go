package mpd

// AbsoluteStart returns the absolute start time of a Period in seconds.
// It is meant for dynamic MPDs where each Period has a start time,
// but it can also be used for the start-over case where a dynamice MPD
// has just turned into a static MPD and still have availabilityStartTime,
// with the same value as in the dynamice case.
func (p *PeriodType) AbsoluteStart(m *MPD) (float64, error) {
	if m.AvailabilityStartTime == "" {
		return 0, ErrASTRequired
	}
	ast, err := m.AvailabilityStartTime.ConvertToSeconds()
	if err != nil {
		return 0, err
	}
	included := false
	for _, period := range m.Periods {
		if period == p {
			included = true
			break
		}
	}
	if !included {
		return 0, ErrPeriodNotFound
	}
	return ast + p.Start.Seconds(), nil
}
