package timestore

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	DiskStore     string        // which directory to store the dump files on disk
	Versions      int           // how many flushes to keep
	FlushPeriod   time.Duration // how often to flush to disk
	FlushVersions int           // how many flushed versions to keep
	CleanupPeriod time.Duration // how often to cleanup flushes
}

func NewConfig() Config {
	var c Config
	c.Versions = 10
	c.FlushPeriod = 1 * time.Minute
	c.CleanupPeriod = 10 * time.Minute
	return c
}

// ensure the configuration is valid
func ValidateConfig(config Config) error {
	err := os.MkdirAll(config.DiskStore, 0700)
	if err != nil {
		return errors.New("Invalid disk store: " + err.Error())
	}
	return nil
}
