package mpd

import "time"

// Ptr returns a pointer to any value
func Ptr[T any](v T) *T {
	return &v
}

// Seconds2DurPtr returns a pointer to a duration given a time in seconds
func Seconds2DurPtr(seconds int) *Duration {
	d := Duration(time.Duration(seconds) * time.Second)
	return &d
}
