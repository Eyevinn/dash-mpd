package mpd

import (
	"math"
	"time"
)

// Ptr returns a pointer to any value
func Ptr[T any](v T) *T {
	return &v
}

// Seconds2DurPtr returns a pointer to a duration given a time in seconds
func Seconds2DurPtr(seconds int) *Duration {
	d := Duration(time.Duration(seconds) * time.Second)
	return &d
}

// Seconds2DurPtrFloat64 returns a pointer to a duration given a float64 time in seconds.
func Seconds2DurPtrFloat64(seconds float64) *Duration {
	us := time.Duration(math.Round(seconds * 1_000_000))
	d := Duration(us * time.Microsecond)
	return &d
}
