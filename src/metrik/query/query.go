package query

import (
	"errors"
	"strconv"
	"time"

	"metrik/projects"

	"github.com/garyburd/redigo/redis"
)

func ParseStartParam(rawStart string) (int64, error) {
	if rawStart != "" {
		s, err := strconv.Atoi(rawStart)
		if err != nil {
			return -1, errors.New("Invalid value for param: start")
		}
		return int64(s), nil
	} else {
		return -1, errors.New("Invalid value for param: start")
	}
}

func ParseEndParam(rawEnd string) (int64, error) {
	if rawEnd != "" {
		if i, err := strconv.Atoi(rawEnd); err == nil {
			return int64(i), nil
		} else {
			return -1, errors.New("Invalid value for param: end")
		}
	} else {
		return time.Now().Unix() * 1000, nil
	}
}

// TODO: better error handling
func Query(r redis.Conn, projectId int, eventName string, start, end int64) ([]string, error) {
	eventsKey := projects.GetEventKey(r, projectId, eventName)

	rawEvents, err := redis.Strings(r.Do("ZRANGEBYSCORE", eventsKey, start, end))
	if err != nil {
		return nil, err
	}

	return rawEvents, nil
}
