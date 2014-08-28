package utils

import (
	"strconv"
	"time"
)

var months = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

func GetYearNo(now time.Time) string {
	return strconv.Itoa(now.Year())
}

// Get the current month number 1-12, 0-padded to 2 digits
func GetMonthNo(now time.Time) string {
	monthNo := strconv.Itoa(months[now.Month().String()])
	if len(monthNo) == 1 {
		monthNo = "0" + monthNo
	}
	return monthNo
}

func GetDayNo(now time.Time) string {
	_, _, dayNo := now.Date()
	return strconv.Itoa(dayNo)
}

func GetHourNo(now time.Time) string {
	return strconv.Itoa(now.Hour())
}

func GetMilliTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
