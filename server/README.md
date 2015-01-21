# Obelisk Data Server

The `server` package implements the obelisk time-series server. This is the persistent data store backing Obelisk's web functionality.

## API
The data server exposes two APIs to the data: a TCP based protocol similar to that of Redis and an HTTP interface that outputs structured JSON.

### TCP

* `declare id path1 path2 ...`: 
* `schema id metric type unit description`: Declare the schema for a metric
* `record id metric time value`: Record a measurement

Example interaction:

    declare a712371 host/h1/service/s1•get service/s1/host/h1•get
    declare a712372 host/h1/service/s1•set service/s1/host/h1•set
    
    schema a712371 counter op 'number of get commands'
    schema a712372 counter op 'number of set commands'
    
    record a712371 2014-05-12 19
    record a712372 2014-05-13 19

## Data Stores
The server is comprised of three types of stores, each backed with their own LevelDB instance. 

### Key-Value
The key-value store enables efficient access for metadata on the schema. It supports `set`/`get`.

### Tag
The tag store enables efficient span queries against trees, it also provides persistent, unique IDs for unique paths (and so provides the keys into the timeseries database). It's primary role is enabling efficient discovery of metrics.

* `worker/<worker>/<metric>`: get metrics for a worker
* `host/<hostname>/<worker>`: get workers by hostname

### Timeseries
The timeseries store is indexed by `uint64`. These ids are provided by the tag store from `worker/<worker>/<metric>` and are guaranteed to be unique for a worker. It provides efficient querying of a particular key over a time range.


## Queries
#### all metrics reported from workers on a host
* query tagstore 