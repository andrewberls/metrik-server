package projects

import (
	"fmt"
	"metrik/utils"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Use a project API key to look up its ID and generate an event
// key, including a GUID
//
// Ex: "projects:abc123:events:signup"
//
func GetEventKey(r redis.Conn, projectKey string, eventName string) string {
	return fmt.Sprintf("projects:%s:events:%s", projectKey, eventName)
}

// Generate a key to track the count of an event for this hour,
// keyed by year / month / day / hour
//
// Ex: "projects:abc123:events:2013-08-22-13"
//
func GetEventCountKey(eventKey string) string {
	now := time.Now()
	yearNo := utils.GetYearNo(now)
	monthNo := utils.GetMonthNo(now)
	dayNo := utils.GetDayNo(now)
	hourNo := utils.GetHourNo(now)
	return fmt.Sprintf("%s:%s-%s-%s-%s",
		eventKey, yearNo, monthNo, dayNo, hourNo)
}

// Generate a key to track the number of events for this project this month,
// keyed by year / month
//
// Ex: "projects:abc123:events:08-2014"
//
func GetProjectEventsCountKey(projectKey string) string {
	now := time.Now()
	monthNo := utils.GetMonthNo(now)
	yearNo := utils.GetYearNo(now)
	return fmt.Sprintf("projects:%s:events:%s-%s",
		projectKey, yearNo, monthNo)
}
