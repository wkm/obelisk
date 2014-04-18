package resp

import (
	"bytes"
	"github.com/bmizerany/assert"
	"reflect"
	"testing"
)

func TestWrites(t *testing.T) {
	testcases := []struct {
		kind reflect.Kind
		val  interface{}
		out  string
	}{
		{reflect.Int, 45, ":45\n\r"},
		{reflect.Int, nil, "$-1\r\n"},

		{reflect.String, "", "$0\n\r\n\r"},
		{reflect.String, "oh hai", "$6\n\roh hai\n\r"},
	}

	for _, tc := range testcases {
		t.Logf("testcase:%q", tc)
		var bb bytes.Buffer

		val := reflect.ValueOf(tc.val)
		_, err := write(&bb, tc.kind, val)
		assert.Equal(t, nil, err)
		assert.Equal(t, tc.out, bb.String())
	}
}
