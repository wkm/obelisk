package resp

import (
	"github.com/bmizerany/assert"
	"reflect"
	"strings"
	"testing"
)

// Validate the fast size compute function
func TestReadSize(t *testing.T) {
	testcases := []struct {
		str string
		sz  int
	}{
		{"\r\n", 0},
		{"1\r\n", 1},
		{"12\r\n", 12},
		{"0123\r\n", 123},
		{"871\r\n", 871},
	}

	for _, tc := range testcases {
		t.Logf("testcase = %v", tc)
		sz, err := readSize(strings.NewReader(tc.str))

		assert.Equal(t, nil, err)
		assert.Equal(t, tc.sz, sz)
	}
}

// Assert that the parsed values matches an arbitrarily typed array
func assertForm(t *testing.T, exp []interface{}, recv []reflect.Value) {
	assert.Equal(t, len(exp), len(recv), "unequal length of form")

	for i := range exp {
		v := reflect.ValueOf(exp[i])
		r := recv[i]

		assert.Equalf(t, v.Kind(), r.Kind(), "exp=%q recv=%q", v, r)
		switch v.Kind() {
		case reflect.String:
			assert.Equal(t, v.String(), r.String())

		default:
			panic("uncomparable type")
		}
	}
}

func TestRead(t *testing.T) {
	testcases := []struct {
		input string
		form  []interface{}
	}{
		{
			"*2\r\n$3\r\nfoo\r\n$4\r\nbraz\r\n",
			[]interface{}{"foo", "braz"},
		},
	}

	for _, tc := range testcases {
		t.Logf("testcase = %v", tc)
		values, err := readRequest(strings.NewReader(tc.input))

		assert.Equal(t, nil, err)
		assertForm(t, tc.form, values)
	}
}
