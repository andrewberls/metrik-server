package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"metrik/projects"
	"metrik/utils"

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
func RangeQuery(r redis.Conn, projectId int, eventName string, start, end int64) ([]string, error) {
	eventsKey := projects.GetEventKey(r, projectId, eventName)

	rawEvents, err := redis.Strings(r.Do("ZRANGEBYSCORE", eventsKey, start, end))
	if err != nil {
		return nil, err
	}

	return rawEvents, nil
}

// Return counts of a given event bucketed by hour
// Ex:
//   {
//     1 => 257,
//     2 => 109
//   }
func HourlyEventCounts(r redis.Conn, projectId int, eventName string, start int64) (string, error) {
	hourlyCounts := make(map[string]int)
	eventKey := projects.GetEventKey(r, projectId, eventName)

	now := time.Now()
	yearNo := utils.GetYearNo(now)
	monthNo := utils.GetMonthNo(now)
	dayNo := utils.GetDayNo(now)

	// TODO: preconstruct keys, MGET
	for i := 1; i < 24; i++ {
		hourNo := strconv.Itoa(i)
		if len(hourNo) == 1 {
			hourNo = "0" + hourNo
		}
		eventKey := fmt.Sprintf("%s:%s-%s-%s-%s",
			eventKey, yearNo, monthNo, dayNo, hourNo)
		count, err := redis.Int(r.Do("GET", eventKey))
		if err != nil {
			// Key missing = 0 events
			hourlyCounts[strconv.Itoa(i)] = 0
		}
		hourlyCounts[strconv.Itoa(i)] = count
	}

	marshalledCounts, err := json.Marshal(hourlyCounts)
	if err != nil {
		return "", err
	}
	return string(marshalledCounts), nil
}

//func HourlyEvents(string, error) {
//        // {
//        //   1 => [{$id=>"a"}, {$id=>"b"}],
//        //   2 => [{$id=>"a"}, {$id=>"b"}],
//        //   3 => [{$id=>"a"}, {$id=>"b"}]
//        // }
//}
