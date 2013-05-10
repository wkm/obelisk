/*
	logging facility for workers
*/
package rlog

import (
	"path/filepath"
	"sync"
)

type Log interface {
	// write a message to the log
	Printf(format string, obj ...interface{})

	// ensure the log's buffer has been synchronized
	Sync()

	// close the log
	Close()
}

// the log proxy encodes a log's configuration along with a
// delegate which actually preforms the logging
type LogProxy struct {
	Category string
	Prefix   string
	Delegate *Log
}

func (p LogProxy) Printf(format string, obj ...interface{}) {
	(*p.Delegate).Printf(format, obj...)
}
func (p LogProxy) Sync() {
	(*p.Delegate).Sync()
}
func (p LogProxy) Close() {
	(*p.Delegate).Close()
}

type Config struct {
	sync.Mutex
	BaseDir *string
	logs    map[string]*LogProxy
}

func NewConfig() *Config {
	c := Config{}
	c.logs = make(map[string]*LogProxy)
	return &c
}

var LogConfig = NewConfig()

// var log = config.Logger("rlog")

// configure the global logger
func ConfigureLogger(c *Config) {
	rlog := c.Logger("rlog")
	rlog.Printf("configuring %v", c)

	c.Lock()
	defer c.Unlock()
	LogConfig = c
}

// given a log proxy, will fill it out
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
		println("received error creating %v: %s", p, err.Error())
	}

	// FIXME threadsafe? we don't lock the logproxy
	p.Delegate = &logger
	return p
}

// get a logger for a named category
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
