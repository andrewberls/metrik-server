package query

import (
	"errors"
	"metrik/projects"
	"strconv"
	"time"

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
		return time.Now().Unix(), nil
	}
}

// TODO: better error handling
func Query(r redis.Conn, apiKey string, eventName string, start, end int64) ([]string, error) {
	//fmt.Println("query.query!", start, end)

	eventsKey, err := projects.GetEventKey(r, apiKey, eventName)
	if err != nil {
		return nil, err
	}

	rawEvents, err := redis.Strings(r.Do("ZRANGEBYSCORE", eventsKey, start, end))
	if err != nil {
		return nil, err
	}

	return rawEvents, nil
}
