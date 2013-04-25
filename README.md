# Obelisk
Obelisk is a management interface for Gocircuit deployments.

## Components
* `/web` -- a web interface 
* `/agent` -- a small worker which is intended to run on every host; provides
  system level monitoring (memory usage, network traffic, disk space, etc.)

* `/rlog` -- a logging utility library for workers that exposes a `remote-log`
  service

Future:

* `/lib` -- the core implementation of obelisk functions
* `/cmd` -- a command line interface that exposes obelisk functionality in a
  programmatic environment
