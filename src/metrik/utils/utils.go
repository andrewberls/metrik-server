package utils

import (
	"strconv"
	"strings"
	"time"
)

const (
	Hello = "world"
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
	monthNo := months[now.Month().String()]
	return Rjust(strconv.Itoa(monthNo), 2, "0")
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

func ToMilliTimestamp(t time.Time) int64 {
	return t.Unix() * 1000
}

// Right-justify <str> and pad to <length> with <padstr>
// Ex:
//   Rjust("8", 2, "0") => "08"
func Rjust(str string, length int, padstr string) string {
	return strings.Repeat(padstr, length-len(str)) + str
}

// TODO - this is a hack
// WE have to combine key and fields into a single slice to splat,
// and need to convert []string ids to []interface{}
func HmgetArgs(key string, fields []string) []interface{} {
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, v := range fields {
		args[i+1] = interface{}(v)
	}
	return args
}
