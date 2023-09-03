package time_utils

import "time"

const (
	Second = 1
	Minute = 60 * Second
	Hour   = 60 * Minute
	Day    = 24 * Hour
)

// DayBeginTime 返回开始的时间戳
func DayBeginTime(sec int64) int64 {
	t := time.Unix(sec, 0)
	return sec - int64(t.Hour())*Hour - int64(t.Minute())*Minute - int64(t.Second())
}

// DayEndTime 返回结束的时间戳
func DayEndTime(sec int64) int64 {
	t := time.Unix(sec, 0)
	return sec + int64(23-t.Hour())*Hour + int64(59-t.Minute())*Minute + int64(59-t.Second())
}
