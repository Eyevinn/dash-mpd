package mpd

import "time"

func StringPtr(v string) *string {
	return &v
}

func IntPtr(v int) *int {
	return &v
}

func Int64Ptr(v int64) *int64 {
	return &v
}

func UintPtr(v uint) *uint {
	return &v
}

func Uint32Ptr(v uint32) *uint32 {
	return &v
}

func Uint64Ptr(v uint64) *uint64 {
	return &v
}

func BoolPtr(v bool) *bool {
	return &v
}

func Float64Ptr(v float64) *float64 {
	return &v
}

// DurPtr is a helper function to generate a pointer to a Duration value.
func DurPtr(v Duration) *Duration {
	return &v
}

// Seconds2DurPtr is a helper function to
func Seconds2DurPtr(seconds int) *Duration {
	d := Duration(time.Duration(seconds) * time.Second)
	return &d
}
