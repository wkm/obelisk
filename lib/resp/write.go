package resp

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

func write(w io.Writer, kind reflect.Kind, val reflect.Value) (nn int, err error) {
	if !val.IsValid() {
		return writeNil(w)
	}

	switch kind {
	case reflect.Int:
		return writeInt(w, val)

	case reflect.String:
		return writeBulkString(w, val)
	}

	return 0, errors.New(fmt.Sprintf("Unsupported kind %q", val.Kind()))
}

func writeInt(w io.Writer, val reflect.Value) (nn int, err error) {
	return fmt.Fprintf(w, ":%d\n\r", val.Int())
}

func writeOk(w io.Writer) (nn int, err error) {
	return fmt.Fprintf(w, "+OK\r\n")
}

func writeSimpleString(w io.Writer, val reflect.Value) (nn int, err error) {
	return fmt.Fprintf(w, "+%s\n\r", val.String())
}

func writeNil(w io.Writer) (nn int, err error) {
	return fmt.Fprintf(w, "$-1\r\n")
}

func writeBulkString(w io.Writer, val reflect.Value) (nn int, err error) {
	str := val.String()
	return fmt.Fprintf(w, "$%d\n\r%s\n\r", len(str), str)
}

func writeArray(w io.Writer, val reflect.Value) (nn int, err error) {
	nn, err = fmt.Fprintf(w, "*%d\r\n", val.Len())
	if err != nil {
		return
	}

	for i := 0; i < val.Len(); i++ {
		i := val.Index(i)
		cnn, err := write(w, val.Kind(), i)
		nn += cnn
		if err != nil {
			return nn, err
		}
	}

	return
}

func writeError(w io.Writer, val *reflect.Value) (nn int, err error) {
	return fmt.Fprintf(w, "-%s\n\r", val.String())
}
