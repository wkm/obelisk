package resp

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/bmizerany/assert"
)

func TestWrites(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		kind reflect.Kind
		val  interface{}
		out  string
	}{
		{reflect.Int, 45, ":45\r\n"},
		{reflect.Int, nil, "$-1\r\n"},

		{reflect.String, "", "$0\r\n\r\n"},
		{reflect.String, "oh hai", "$6\r\noh hai\r\n"},
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
