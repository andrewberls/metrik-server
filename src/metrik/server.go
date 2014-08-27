package main

import (
	"net/http"
	"strings"

	"metrik/ingest"
	"metrik/projects"
	"metrik/query"

	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func main() {
	m := martini.Classic()
	r, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	m.Post("/v1/events", binding.Bind(ingest.EventParams{}),
		func(eventParams ingest.EventParams) (int, string) {
			err := ingest.IngestEvent(r, eventParams)
			if err != nil {
				return 400, err.Error() // TODO: better status code handling
			}
			return 201, ""
		})

	m.Get("/v1/events", func(req *http.Request) (int, string) {
		params := req.URL.Query()

		apiKey := params.Get("api_key")
		if apiKey == "" {
			return 400, "Invalid API key"
		}

		eventName := params.Get("event")
		if eventName == "" {
			return 400, "Invalid event name"
		}

		start, err := query.ParseStartParam(params.Get("start"))
		if err != nil {
			return 400, err.Error()
		}

		end, err := query.ParseEndParam(params.Get("end"))
		if err != nil {
			return 400, err.Error()
		}

		projectId, err := projects.GetProjectId(r, apiKey)
		if err != nil {
			return 400, "Invalid API key"
		}

		events, err := query.Query(r, projectId, eventName, start, end)
		if err != nil {
			return 400, err.Error() // TODO: better status code handling
		}
		json := "[" + strings.Join(events, ",") + "]"

		return 200, json
	})

	m.Run()
}
