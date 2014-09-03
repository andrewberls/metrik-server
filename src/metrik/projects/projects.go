package projects

import (
	"fmt"
	"time"

	"metrik/utils"
)

// Hash key for id => json lookups
func GetEventsKey(projectKey string) string {
	return fmt.Sprintf("projects:%s:events", projectKey)
}

// TODO: need name here!
// Zset key for timestamp => id lookups
func GetEventTimesKey(projectKey string, eventName string) string {
	return fmt.Sprintf("projects:%s:event_times:%s", projectKey, eventName)
}

// List key for name => [ids] lookups
func GetEventNameKey(projectKey string, eventName string) string {
	return fmt.Sprintf("projects:%s:event:%s", projectKey, eventName)
}

// TODO: this probably needs an event name ?
//
// Generate a key to track the count of an event for this hour,
// keyed by year / month / day / hour
//
// Ex: "projects:abc123:events:2013-08-22-13"
//
//func GetEventCountKey(project string) string {
//        eventKey := GetEvents
//        now := time.Now()
//        yearNo := utils.GetYearNo(now)
//        monthNo := utils.GetMonthNo(now)
//        dayNo := utils.GetDayNo(now)
//        hourNo := utils.GetHourNo(now)
//        return fmt.Sprintf("%s:%s-%s-%s-%s",
//                eventKey, yearNo, monthNo, dayNo, hourNo)
//}

// Integer key to count project event quota, keyed by year/month
//
// Ex: "projects:abc123:events:08-2014"
//
func GetProjectEventsCountKey(projectKey string) string {
	now := time.Now()
	monthNo := utils.GetMonthNo(now)
	yearNo := utils.GetYearNo(now)
	return fmt.Sprintf("projects:%s:event_count:%s-%s",
		projectKey, yearNo, monthNo)
}
