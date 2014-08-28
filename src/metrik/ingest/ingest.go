package ingest

import (
	"encoding/json"
	"errors"

	"metrik/projects"
	"metrik/utils"

	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
)

type EventParams struct {
	Name       string `form:"event" binding:"required"`
	ApiKey     string `form:"api_key" binding:"required"`
	ProjectKey string `form:"project_key" binding:"required"`
	Properties string `form:"properties" binding:"required"`
}

// Convert a raw property JSON string into a map with the correct format of:
//   {
//      "$id" => "9ba23839-82aa-4efe-bc17-2f3ca424cb8f",
//      "$name" => "users.signup"
//      "$timestamp" => 192837387876,
//      "$properties" => "{'hello':'world'}",
//   }
func formatEvent(timestamp int64, eventParams EventParams) (map[string]interface{}, error) {
	event := make(map[string]interface{})

	var properties map[string]interface{}
	if err := json.Unmarshal([]byte(eventParams.Properties), &properties); err != nil {
		return nil, err // Invalid client-supplied JSOn
	}

	event["$id"] = uuid.New()
	event["$name"] = eventParams.Name
	event["$timestamp"] = timestamp
	event["$properties"] = properties

	return event, nil
}

// Main entry point to ingest an event given raw parameters
// Transforms params into standard form hash, stores event in redis zset
// using millitimestamp as score, counts the event for the project's monthly
// quota, and counts the event for this hour
func IngestEvent(r redis.Conn, eventParams EventParams) error {
	projectKey := eventParams.ProjectKey
	eventKey := projects.GetEventKey(r, projectKey, eventParams.Name)

	milliTimestamp := utils.GetMilliTimestamp()
	event, err := formatEvent(milliTimestamp, eventParams)
	if err != nil {
		return errors.New("Invalid property JSON")
	}

	marshalledEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	score := milliTimestamp

	r.Send("MULTI")
	r.Send("ZADD", eventKey, score, string(marshalledEvent))
	r.Send("INCR", projects.GetEventCountKey(eventKey))
	r.Send("INCR", projects.GetProjectEventsCountKey(projectKey))
	r.Do("EXEC")

	return nil
}
