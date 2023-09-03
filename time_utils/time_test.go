package time_utils

import (
	"testing"
	"time"
)

func TestBeginTime(t *testing.T) {
	str := time.Unix(DayBeginTime(time.Now().Unix()), 0).Format("2006-01-02 15:04:05")
	t.Log(str)
}
func TestEndTime(t *testing.T) {
	str := time.Unix(DayEndTime(time.Now().Unix()), 0).Format("2006-01-02 15:04:05")
	t.Log(str)
}
