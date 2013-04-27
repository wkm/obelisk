# Obelisk
Obelisk is a management interface for Gocircuit deployments.

## Components
* `/web` -- a web interface to obelisk
* `/agent` -- a small worker which is intended to run on every host; provides
  system level monitoring (memory usage, network traffic, disk space, etc.)

* `/server` -- the obelisk server

* `/cmd` -- obelisk utilities
    * `logs`: utility for remotely reading logs from workers

* `/lib` -- backbone libraries for obelisk
    * **`rinst`**: a library for instrumenting code with metrics
    * **`rlog`**: a logging utility library for workers that exposes a `remote-log` service

    * **`storekv`**: a simple key-value database
    * **`storetag`**: a simple tags database
    * **`storetime`**: a simple timestore database
    * **`persist`**: a utility library for creating simple persisted databases

    * `rconfig`: [experimental] remote configuration stuffs
    * `rinstreporter`: [stupid] just... stupid


## Consuming Measurements

```go
ch := make(chan Vars, 1000)

for {
	time.Sleep(time.Second)
	for r := fetchVars() {
		r.write(ch)
	}
}

for {
	w <- ch
	vena.Write(w)
}
```