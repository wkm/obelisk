// remote configuration facility for workers
package rconfig

import (
	"circuit/use/circuit"
	"sync"
)

const ServiceName = "remote-config"

// a simple configuration object
type RConfig struct {
	sync.Mutex
	Settings map[string]string
}

var Config = new(RConfig)

func init() {
	circuit.RegisterValue(&RConfig{})
}

// change the value of a setting (threadsafe)
func (r *RConfig) Set(key, value string) {
	r.Lock()
	defer r.Unlock()

	r.Settings[key] = value
}

// get the value of a setting (threadsafe)
func (r *RConfig) Get(key string) string {
	r.Lock()
	defer r.Unlock()

	return r.Settings[key]
}

// get a copy of all values set (threadsafe)
func (r *RConfig) GetAll(key string) map[string]string {
	r.Lock()
	defer r.Unlock()

	return r.Settings
}
