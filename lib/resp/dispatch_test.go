package resp

import (
	"testing"
)

type keyval struct {
	values map[string]string
}

func (kv *keyval) Get(key string) string {
	return kv.values[key]
}

func (kv *keyval) Put(key, value string) {
	kv.values[key] = value
}

func TestDispatch(t *testing.T) {
	k := keyval{}
	d, err := newDispatch(&k)
	if err != nil {
		t.Log(err.Error())
	}

	t.Logf("%v", d)
}
