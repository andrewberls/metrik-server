package query

import (
	"encoding/json"
	"errors"
	"metrik/projects"
	"metrik/utils"
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
		return time.Now().Unix() * 1000, nil
	}
}

// Return all occurrences of <eventName> in the timespan defined by start and end
func RangeQuery(r redis.Conn, projectKey string, eventName string, start, end int64) ([]string, error) {
	eventTimesKey := projects.GetEventTimesKey(projectKey, eventName)

	ids, err := redis.Strings(r.Do("ZRANGEBYSCORE", eventTimesKey, start, end))
	if err != nil {
		return nil, err
	}

	// TODO: this is lame, need to combine key with ids for variadic args
	var args []interface{}
	args = append(args, projects.GetEventsKey(projectKey))
	for _, id := range ids {
		args = append(args, id)
	}

	rawEvents, err := redis.Strings(r.Do("HMGET", args...))
	if err != nil {
		return nil, err
	}

	return rawEvents, nil
}

// TODO: could just count values if HourlyEvents returned map instead of raw string
//
// Return counts of a given event bucketed by hour (as a marshalled JSON string)
// Ex format:
//   {
//     1 => 257,
//     2 => 109,
//     ...
//   }
func HourlyEventCounts(r redis.Conn, projectKey string, eventName string, start int64) (string, error) {
	return "", nil
	//hourlyCounts := make(map[string]int)
	//eventKey := projects.GetEventKey(r, projectKey, eventName)

	//now := time.Now()
	//yearNo := utils.GetYearNo(now)
	//monthNo := utils.GetMonthNo(now)
	//dayNo := utils.GetDayNo(now)

	//// Preconstruct hourly event keys and MGET
	//hourKeys := make([]interface{}, 24) // r.Do needs []interface{}, not []string
	//for i := 1; i < 24; i++ {
	//        hourNo := utils.Rjust(strconv.Itoa(i), 2, "0")
	//        eventKey := fmt.Sprintf("%s:%s-%s-%s-%s",
	//                eventKey, yearNo, monthNo, dayNo, hourNo) // TODO: all keygen should go through projects
	//        hourKeys[i] = eventKey
	//}

	//eventCounts, err := redis.Strings(r.Do("MGET", hourKeys...))
	//for idx, strCount := range eventCounts {
	//        if strCount == "" {
	//                hourlyCounts[strconv.Itoa(idx)] = 0
	//        } else {
	//                intCount, _ := strconv.Atoi(strCount)
	//                hourlyCounts[strconv.Itoa(idx)] = intCount
	//        }
	//}

	//marshalledCounts, err := json.Marshal(hourlyCounts)
	//if err != nil {
	//        return "", err
	//}
	//return string(marshalledCounts), nil
}

// TODO
// Ex format:
//   {
//     1 => [{$id=>"a"}, {$id=>"b"}],
//     2 => [{$id=>"a"}, {$id=>"b"}],
//     3 => [{$id=>"a"}, {$id=>"b"}],
//     ...
//   }
//
func HourlyEvents(r redis.Conn, projectKey string, eventName string, start time.Time) (string, error) {
	hourlyEvents := make(map[string][]string)

	eventsKey := projects.GetEventsKey(projectKey)                    // Hash id -> <json>
	eventTimesKey := projects.GetEventTimesKey(projectKey, eventName) // Zset timestamp -> id

	midnight := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	lower := midnight
	upper := lower.Add(time.Hour).Add(-time.Second) // TODO: inclusivity

	// TODO: parallelize lookups
	var ids []string
	for i := 0; i < 24; i++ {
		reply, err := redis.Values(r.Do("ZRANGEBYSCORE",
			eventTimesKey,
			utils.ToMilliTimestamp(lower),
			utils.ToMilliTimestamp(upper)))
		if err != nil {
			panic(err)
		}
		if err := redis.ScanSlice(reply, &ids); err != nil {
			panic(err)
		}

		if len(ids) == 0 {
			hourlyEvents[strconv.Itoa(i)] = []string{}
		} else {
			args := utils.HmgetArgs(eventsKey, ids)
			rawEvents, err := redis.Strings(r.Do("HMGET", args...))
			if err != nil {
				panic(err)
			}
			hourlyEvents[strconv.Itoa(i)] = rawEvents
		}

		lower = upper.Add(time.Second)
		upper = lower.Add(time.Hour).Add(-time.Second)
	}

	marshalledEvents, err := json.Marshal(hourlyEvents)
	if err != nil {
		return "", err
	}
	return string(marshalledEvents), nil
}
