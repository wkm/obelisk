package resp

import (
	"errors"
	"testing"
)

type keyval struct {
	values map[string]string
}

func newKeyval() (kv *keyval) {
	kv = &keyval{}
	kv.values = make(map[string]string)
	return
}

func (kv *keyval) Get(key string) string {
	return kv.values[key]
}

func (kv *keyval) MGet(k1, k2 string) []string {
	return []string{
		kv.values[k1],
		kv.values[k2],
	}
}

func (kv *keyval) Put(key, value string) {
	kv.values[key] = value
}

func (kv *keyval) Err() error {
	return errors.New("Testing error")
}

func TestStringDispatch(t *testing.T) {
	k := newKeyval()
	d, err := newDispatch(k)
	if err != nil {
		t.Log(err.Error())
	}

	testcases := []struct{ in, out string }{
		{"get hi", "$0\r\n\r\n"},
		{"put hi there", "+OK\r\n"},
		{"get hi", "$5\r\nthere\r\n"},
		{"err", "-Testing error\r\n"},
		{"MGET hi dog", "*2\r\n$5\r\nthere\r\n$0\r\n\r\n"},
	}

	for _, tc := range testcases {
		t.Logf("testcase: %q", tc)

		res, err := d.Call(tc.in)
		if err != nil {
			t.Errorf("unexpected error %q", err)
		}

		if res != tc.out {
			t.Errorf("expected %q; received %q", tc.out, res)
		}
	}
}
