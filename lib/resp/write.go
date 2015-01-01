package resp

import (
	"fmt"
	"io"
	"reflect"
)

type Response interface {
	Write(w io.Writer)
}

func write(w io.Writer, kind reflect.Kind, val reflect.Value) (nn int, err error) {
	if !val.IsValid() {
		return writeNil(w)
	}

	switch kind {
	case reflect.Int:
		return writeInt(w, val.Int())

	case reflect.String:
		return writeBulkString(w, val.String())

	case reflect.Interface:
		errorMeth := val.MethodByName("Error")
		if errorMeth.Type().NumIn() == 0 && errorMeth.Type().NumOut() == 1 {
			out := errorMeth.Call([]reflect.Value{})
			return writeError(w, out[0].String())
		}

	case reflect.Slice, reflect.Array:
		return writeArray(w, val)
	}

	return 0, fmt.Errorf("Unsupported kind %q and type %q", val.Kind(), val.Type())
}

func writeInt(w io.Writer, val int64) (nn int, err error) {
	return fmt.Fprintf(w, ":%d\r\n", val)
}

func writeOk(w io.Writer) (nn int, err error) {
	return fmt.Fprintf(w, "+OK\r\n")
}

func writeSimpleString(w io.Writer, val string) (nn int, err error) {
	return fmt.Fprintf(w, "+%s\r\n", val)
}

func writeNil(w io.Writer) (nn int, err error) {
	return fmt.Fprintf(w, "$-1\r\n")
}

func writeBulkString(w io.Writer, str string) (nn int, err error) {
	return fmt.Fprintf(w, "$%d\r\n%s\r\n", len(str), str)
}

func writeArray(w io.Writer, val reflect.Value) (nn int, err error) {
	nn, err = fmt.Fprintf(w, "*%d\r\n", val.Len())
	if err != nil {
		return
	}

	for i := 0; i < val.Len(); i++ {
		i := val.Index(i)
		cnn, err := write(w, i.Kind(), i)
		nn += cnn
		if err != nil {
			return nn, err
		}
	}

	return
}

func writeError(w io.Writer, val string) (nn int, err error) {
	return fmt.Fprintf(w, "-Error: %s\r\n", val)
}
