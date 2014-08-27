package projects

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

// TODO: this is lame
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

// Look up a project's ID using its API key
// TODO: this lookup is lame, use project key ? or specify api key + project id ?
//
// Format: "projects:<api_key>"
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
// Format: "projects:<id>:events:<name>"
//
func GetEventKey(r redis.Conn, projectId int, name string) string {
	return fmt.Sprintf("projects:%d:events:%s", projectId, name)
}

// Generate a key to track the number of events for this project this month
//
// Ex: "projects:2:events:08-2014"
//
func GetProjectEventCountKey(projectId int) string {
	now := time.Now()
	monthNo := strconv.Itoa(months[now.Month().String()])
	if len(monthNo) == 1 {
		monthNo = "0" + monthNo
	}
	year := strconv.Itoa(now.Year())
	return fmt.Sprintf("projects:%d:events:%s-%s", projectId, monthNo, year)
}
