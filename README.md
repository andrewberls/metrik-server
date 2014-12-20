# Metrik

Metrik is an event intake and query server. It was written primarily as an exercise
in learning the [Go](https://golang.org/) language, and is not by any means complete nor production-ready.
It uses [Martini](https://github.com/go-martini/martini)
for the web layer and [Redis](http://redis.io/) for storing event data.

## Getting Started (Mac)

```
$ go get ./...
$ ./bin/gox -osarch="darwin/amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" metrik
$ redis-server
$ ./bin/metrik_darwin_amd64
```

### Create an event
```
$ curl -X POST -d 'event=login&project_key=abc&properties={"hello": "world"}' http://localhost:3000/v1/events/
// 201 Created

$ curl -X POST -d 'event=login&project_key=abc&properties={"hi": "again"}' http://localhost:3000/v1/events/
// 201 Created
```

### Querying events
```
// start param in Unix time, end param is optional
$ curl -X GET http://localhost:3000/v1/events?project_key=abc&event=login&start=1418976000
// [{"$id":"8164a424-0a4f-44db-8279-ae298a05b12a","$name":"login","$properties":{"hello":"world"},"$timestamp":1419065456051},
    {"$id":"97a6b929-f966-4354-a2aa-5e16075c7d4b","$name":"login","$properties":{"hi":"again"},"$timestamp":1419065718486}]

$ curl -X GET http://localhost:3000/v1/hourly_events?project_key=abc&event=login&start=1418976000
// {"0" ["{\"$id\":\"8164a424-0a4f-44db-8279-ae298a05b12a\", ...}"]
    "1" [...]
    ...
    "23" [...]}
```


