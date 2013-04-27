package timestore

import (
	"time"
)

type Config struct {
	InitialSize int           // the initial timeline capacity
	DiskStore   string        // which directory to store the dump files on disk
	Versions    int           // how many flushes to keep
	FlushPeriod time.Duration // how often to flush to disk
}
