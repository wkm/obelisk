package resp

import (
	"bytes"
	"reflect"
	"testing"
)

type container struct{}

func (container) NullOutput() {}

func (container) SingleOutput() int {
	return 0
}

func (container) TooManyOutput() (int, int, int) {
	return 0, 0, 0
}

func (container) NonErrorSecondOutput() (int, int) {
	return 0, 0
}

func (container) SingleErrorOutput() error {
	return nil
}

func (container) WithErrorOutput() (int, error) {
	return 0, nil
}

func (container) InvalidOutputType() bytes.Buffer {
	return bytes.Buffer{}
}

func (container) InvalidInputType(bb bytes.Buffer) {}

func (container) StringArrayInput([]string) {}

func TestMethodSignature(t *testing.T) {
	t.Parallel()

	ty := reflect.TypeOf(container{})

	testcases := []struct {
		methodName string
		expected   error
	}{
		{"NullOutput", nil},
		{"SingleOutput", nil},
		{"SingleErrorOutput", nil},
		{"TooManyOutput", ErrInvalidOutputSize},
		{"NonErrorSecondOutput", ErrNonErrorSecondOutput},
		{"WithErrorOutput", nil},
		{"InvalidInputType", ErrUnsupportedInputKind},
		{"InvalidOutputType", ErrUnsupportedOutputKind},
		{"StringArrayInput", nil},
	}

	for _, tc := range testcases {
		t.Logf("Testing: %q", tc)
		m, ok := ty.MethodByName(tc.methodName)
		if !ok {
			t.Errorf("No method with name %q", tc.methodName)
		}

		err := validateMethod(&m)
		if err != tc.expected {
			t.Errorf("Expected err=%q; got=%q", tc.expected, err)
		}
	}
}
