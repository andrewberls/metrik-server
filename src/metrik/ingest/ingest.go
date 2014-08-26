package ingest

import (
	"encoding/json"
	"errors"
	"time"

	"metrik/projects"

	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
)

type EventParams struct {
	Name       string `form:"event" binding:"required"`
	ApiKey     string `form:"api_key" binding:"required"`
	Properties string `form:"properties" binding:"required"`
}

// Decode client-supplied JSON into a map
func unmarshalProperties(rawProperties string) (map[string]interface{}, error) {
	var properties map[string]interface{}
	if err := json.Unmarshal([]byte(rawProperties), &properties); err != nil {
		return nil, err
	}
	return properties, nil
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
	unmarshalledProperties, err := unmarshalProperties(eventParams.Properties)
	if err != nil {
		return nil, err // Invalid client-supplied JSON
	}
	event["$id"] = uuid.New()
	event["$name"] = eventParams.Name
	event["$timestamp"] = timestamp
	event["$properties"] = unmarshalledProperties
	return event, nil
}

// Marshal a map describing an event into JSON
func marshalEvent(event map[string]interface{}) (string, error) {
	marshalledProperties, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	return string(marshalledProperties), nil
}

func IngestEvent(r redis.Conn, eventParams EventParams) error {
	eventKey, err := projects.GetEventKey(r, eventParams.ApiKey, eventParams.Name)
	if err != nil {
		// Couldn't find API key to build event key
		return errors.New("Invalid API key")
	}

	milliTimestamp := time.Now().UnixNano() / 1e6
	event, err := formatEvent(milliTimestamp, eventParams)
	if err != nil {
		return errors.New("Invalid property JSON")
	}

	marshalledEvent, err := marshalEvent(event)
	if err != nil {
		return err
	}

	score := milliTimestamp
	r.Do("ZADD", eventKey, score, marshalledEvent)
	return nil
}
