package resp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestWriteInt(t *testing.T) {
	var bb bytes.Buffer
	val := reflect.ValueOf(45)

	nn, err := writeInt(&bb, &val)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	if bb.String() != ":45\n\r" {
		t.Errorf("wrong string %s", bb.String())
	}

	if nn != 5 {
		t.Errorf("wrong length written: %v", nn)
	}
}

func TestWriteString(t *testing.T) {
	var bb bytes.Buffer
	val := reflect.ValueOf("oh hai")

	nn, err := writeString(&bb, &val)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	if bb.String() != "+oh hai\n\r" {
		t.Errorf("wrong string %s", bb.String())
	}

	if nn != 9 {
		t.Errorf("wrong length written: %v", nn)
	}
}
