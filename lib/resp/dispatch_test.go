package resp

import (
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

func (kv *keyval) Put(key, value string) {
	kv.values[key] = value
}

func TestDispatch(t *testing.T) {
	k := newKeyval()
	d, err := newDispatch(k)
	if err != nil {
		t.Log(err.Error())
	}

	testcases := []struct{ in, out string }{
		{"get hi", "$0\n\r\n\r"},
		{"put hi there", "+OK\r\n"},
		{"get hi", "$5\n\rthere\n\r"},
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
