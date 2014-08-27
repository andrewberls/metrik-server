package projects

import (
	"fmt"
	"metrik/utils"
	"time"

	"github.com/garyburd/redigo/redis"
)

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
func GetEventKey(r redis.Conn, projectId int, eventName string) string {
	return fmt.Sprintf("projects:%d:events:%s", projectId, eventName)
}

// Generate a key to track the count of an event for this hour,
// keyed by year / month / day / hour
//
// Ex: "projects:2:events:2013-08-22-13"
//
func GetEventCountKey(eventKey string) string {
	now := time.Now()
	yearNo := utils.GetYearNo(now)
	monthNo := utils.GetMonthNo(now)
	dayNo := utils.GetDayNo(now)
	hourNo := utils.GetHourNo(now)
	return fmt.Sprintf("%s:%s-%s-%s-%s", eventKey, yearNo, monthNo, dayNo, hourNo)
}

// Generate a key to track the number of events for this project this month,
// keyed by year / month
//
// Ex: "projects:2:events:08-2014"
//
func GetProjectEventsCountKey(projectId int) string {
	now := time.Now()
	monthNo := utils.GetMonthNo(now)
	yearNo := utils.GetYearNo(now)
	return fmt.Sprintf("projects:%d:events:%s-%s",
		projectId, yearNo, monthNo)
}
