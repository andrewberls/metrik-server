package projects

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
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

func getYearNo(now time.Time) string {
	return strconv.Itoa(now.Year())
}

// Get the current month number 1-12, 0-padded to 2 digits
func getMonthNo(now time.Time) string {
	monthNo := strconv.Itoa(months[now.Month().String()])
	if len(monthNo) == 1 {
		monthNo = "0" + monthNo
	}
	return monthNo
}

func getDayNo(now time.Time) string {
	_, _, dayNo := now.Date()
	return strconv.Itoa(dayNo)
}

func getHourNo(now time.Time) string {
	return strconv.Itoa(now.Hour())
}

// Look up a project's ID using its API key
// TODO: use project key ? or client specify api key + project id ?
//
// Ex: "projects:abcd-123"
//
func GetProjectId(r redis.Conn, apiKey string) (int, error) {
	projectKey := fmt.Sprintf("projects:%s", apiKey)
	projectId, err := redis.Int(r.Do("GET", projectKey))
	if err != nil {
		return -1, err
	}
	return projectId, nil
}

// Use a project API key to look up its ID and generate an event
// key, including a GUID
//
// Ex: "projects:2:events:signup"
//
func GetEventKey(r redis.Conn, projectId int, name string) string {
	return fmt.Sprintf("projects:%d:events:%s", projectId, name)
}

// Generate a key to track the count of an event for this hour,
// keyed by year / month / day / hour
//
// Ex: "projects:2:events:2013-08-22-13"
//
func GetEventCountKey(eventKey string) string {
	now := time.Now()
	yearNo := getYearNo(now)
	monthNo := getMonthNo(now)
	dayNo := getDayNo(now)
	hourNo := getHourNo(now)
	return fmt.Sprintf("%s:%s-%s-%s-%s", eventKey, yearNo, monthNo, dayNo, hourNo)
}

// Generate a key to track the number of events for this project this month,
// keyed by year / month
//
// Ex: "projects:2:events:08-2014"
//
func GetProjectEventsCountKey(projectId int) string {
	now := time.Now()
	monthNo := getMonthNo(now)
	yearNo := getYearNo(now)
	return fmt.Sprintf("projects:%d:events:%s-%s",
		projectId, yearNo, monthNo)
}
