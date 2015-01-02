# Obelisk

[![Build Status](https://travis-ci.org/wkm/obelisk.svg?branch=master)](https://travis-ci.org/wkm/obelisk)

Obelisk is a plug-and-play, end-to-end service for measuring, storing,
querying, and visualizing general time-series data for Go programs,
and especially Go programs which run on GoCircuit.

## Service Organization
Obelisk is divided into the following component services:

### `obelisk-server`
The data and query server for Obelisk.

### `obelisk-agent`
A minor utility that use sigar to collect instrumentation data on hosts (CPU usage, memory, network, etc.)

### `obelisk-store`
The store service 

### `obelisk-web`
The web service is (mostly) stateless and can be horizontally scaled to run on as many machines as desired.

### `obelisk-proxy` (future)
The proxy service is a simple (but sufficient) HTTP proxy to the `obelisk-web` and `obelisk-aggregator` services.

### `obelisk-aggregator` (future)
The aggregator service acts as a front end for receiving and aggregating metrics

### `obelisk-broadcast` (future)





## Code Organization
* `/web` a web interface to obelisk
* `/obelisk-cli` a command line utility for 
* `/obelisk-agent` a small worker which is intended to run on every host; provides system level monitoring (memory usage, network traffic, disk space, etc.)
* `/server` -- the obelisk server
* `/cmd` -- obelisk utilities
    * `logs`: utility for remotely reading logs from workers
* `/lib` -- backbone libraries for obelisk
    * **`streamhist`**: a library implementing 
    * **`rinst`**: a library for instrumenting code with metrics
    * **`rlog`**: a logging utility library for workers that exposes a `remote-log` service

    * **`storekv`**: a simple key-value database
    * **`storetag`**: a simple tags database
    * **`storetime`**: a simple timestore database
    * **`persist`**: a utility library for creating simple persisted databases

    * `rconfig`: [experimental] remote configuration stuffs
    * `rinstreporter`: [stupid] just... stupid

## Development
Run `prepare.sh` to setup your local environment for development.
