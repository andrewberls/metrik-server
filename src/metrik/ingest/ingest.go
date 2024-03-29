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
		return nil, err // Invalid client-supplied JSON
	}

	event["$id"] = uuid.New()
	event["$name"] = eventParams.Name
	event["$timestamp"] = timestamp
	event["$properties"] = properties

	return event, nil
}

func IngestEvent(r redis.Conn, eventParams EventParams) error {
	projectKey := eventParams.ProjectKey

	milliTimestamp := utils.GetMilliTimestamp()
	event, err := formatEvent(milliTimestamp, eventParams)
	if err != nil {
		return errors.New("Invalid property JSON")
	}

	marshalledEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	eventId := event["$id"]
	eventName := event["$name"].(string)

	r.Send("MULTI")

	// Hash of id => json (ID lookup)
	r.Send("HSET", projects.GetEventsKey(projectKey), eventId, marshalledEvent)

	// Zset of timestamp => id (Time lookup)
	r.Send("ZADD", projects.GetEventTimesKey(projectKey, eventName), milliTimestamp, eventId)

	// KV of name => [ids] (Name lookup)
	r.Send("RPUSH", projects.GetEventNameKey(projectKey, eventName), eventId)

	// Hourly counter for this event
	// TODO - necessary ?
	//r.Send("INCR", projects.GetEventCountKey(projectKey, eventName))

	// Count event towards project quota
	r.Send("INCR", projects.GetProjectEventsCountKey(projectKey))

	// TODO:
	r.Do("EXEC")

	return nil
}
