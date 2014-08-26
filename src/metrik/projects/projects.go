package projects

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// Look up a project's ID using its API key
func getProjectId(r redis.Conn, apiKey string) (int, error) {
	projectKey := fmt.Sprintf("projects:%s", apiKey)
	projectId, err := redis.Int(r.Do("GET", projectKey))
	if err != nil {
		return -1, err
	}
	return projectId, nil
}

// Use a project API key to look up its ID and generate an event
// key, including a GUID
func GetEventKey(r redis.Conn, apiKey string, name string) (string, error) {
	projectId, err := getProjectId(r, apiKey)
	if err != nil {
		return "", err
	}

	eventKey := fmt.Sprintf("projects:%d:events:%s", projectId, name)
	return eventKey, nil
}
