# Obelisk Data Server

The `server` package implements the obelisk time-series server. This backs

## Data Stores
The server uses three data stores:

### Key-Value
The key value store contains the human-readable data on the schema.

* `worker/<worker>/<metric>`: metric schema object (`unit` and `description`)

### Tag
The tag store is used in two ways: it provides unique IDs for unique paths (and so provides the keys into the timeseries database) and to enable efficient discovery of metrics.

* `worker/<worker>/<metric>`: get metrics for a worker
* `host/<hostname>/<worker>`: get workers by hostname

The two queries are `createNode`

### Timeseries
The timeseries store is indexed by `uint64`. These ids are provided by the tag store from `worker/<worker>/<metric>` and are guaranteed to be unique for a worker.


## Queries
#### all metrics reported from workers on a host
* query tagstore 