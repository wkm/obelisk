// Package rlog provides logging facilities for workers
package rlog

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Log represents a type which
type Log interface {
	// write a message to the log
	Printf(format string, obj ...interface{})

	// ensure the log's buffer has been synchronized
	Sync()

	// close the log
	Close()
}

// LogProxy encodes a log's configuration along with a delegate
// which actually preforms the logging
type LogProxy struct {
	Category string
	Prefix   string
	Delegate *Log
}

// Printf writes a logline into the delegate.
func (p LogProxy) Printf(format string, obj ...interface{}) {
	(*p.Delegate).Printf(format, obj...)
}

// Sync ensures the delegate sync'd.
func (p LogProxy) Sync() {
	(*p.Delegate).Sync()
}

// Close ensures the delegate is closed.
func (p LogProxy) Close() {
	(*p.Delegate).Close()
}

// Config stores configuration for multiple logging categories.
type Config struct {
	sync.Mutex
	BaseDir *string
	logs    map[string]*LogProxy
}

// NewConfig allocates a new logging configuration.
func NewConfig() *Config {
	c := Config{}
	c.logs = make(map[string]*LogProxy)
	return &c
}

// LogConfig
var LogConfig = NewConfig()

// ConfigureLogger configures the global, and default logger
func ConfigureLogger(c *Config) {
	rlog := c.Logger("rlog")
	rlog.Printf("configuring %v", c)

	c.Lock()
	defer c.Unlock()
	LogConfig = c
}

// Create does ... ? what?
func (c *Config) Create(p *LogProxy) Log {
	c.Lock()
	defer c.Unlock()

	if p.Delegate != nil {
		p.Close()
	}

	var logger Log
	var err error
	if c.BaseDir == nil {
		logger, err = StdoutLog{p.Category}, nil
	} else {
		logger, err = NewFileLog(filepath.Join(*c.BaseDir), p.Category, p.Category)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "received error creating %v: %s\n", p, err.Error())
	}

	// FIXME threadsafe? we don't lock the logproxy
	p.Delegate = &logger
	return p
}

// Logger gives a logger for a named category
func (c *Config) Logger(ctg string) Log {
	c.Lock()
	proxy, ok := c.logs[ctg]
	if !ok {
		proxy = &LogProxy{ctg, ctg, nil}
		c.logs[ctg] = proxy
	}

	log := c.logs[ctg]
	c.Unlock()

	c.Create(log)
	return proxy
}
