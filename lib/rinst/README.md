# rinst
A package for instrumenting your application.


## Instruments
### Counter
A `Counter` tracks how many times something has occured. It can only increase. Generally counters become interesting when their change is measured over time. (eg. how many requests have been served)

### Value
A `Value` is a container that exposes how many times a value has changed as well as its current state. (eg. the hostname or number of active workers)

### XXXX Dimension
A `Dimension` is a named set of instruments which share a common property (including other dimensions)

### XXXX Gauge
A `Gauge` is a function which can be executed to get the instantaneous state of the system. (eg. the number of open files)

### XXXX Ratio
A `Ratio` tracks the current total value of an instrument, as well as how the value is allocated across disjoint sets. (eg. the memory available on a system)

### XXXX Metric
A `Metric` tracks a value and it's distribution over time. (eg. the execution time for a request)

### XXXX Error
An `Error` instrument is a kind of counter