# Dump Endomondo Data

```
$ ENDOMONDO_EMAIL=xxx ENDOMONDO_PASSWORD=yyy dumpomondo

== Dumping data for user 111
== Dumping data to JSON (dump/111/raw)
Dumping workouts for user 111, limit 100, offset 0
Dumping workouts for user 111, limit 100, offset 100
...

== Dumping data to GPX (dump/111/gpx)
2020/11/10 17:28:02 Dumping workout 10001 for user 111
2020/11/10 17:28:04 Dumping workout 10002 for user 111
2020/11/10 17:28:05 Dumping workout 10003 for user 111
...
```

## Detailed usage

```
$ dumpomondo -help

Usage of dumpomondo:
  -gpx-sports string
    	sport IDs to be dumped to GPX (default "0,1,2,16,18")
  -skip-gpx
    	skip dumping workout data to GPX
  -skip-json
    	skip dumping workout data to JSON
```

List of supported sports is available in [sports file](sports.go). The sports that are dumped to GPX by default are:
- running,
- cycling (two types),
- hiking,
- walking

## Installation 

```
go get github.com/slomek/dumpomondo
```