// remote configuration facility for workers
package rconfig

import (
	"sync"
)

// a simple configuration object
type RConfig struct {
	sync.Mutex
	Settings map[string]string
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
