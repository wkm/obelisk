package resp

import (
	"fmt"
	"io"
	"reflect"
)

func writeInt(w io.Writer, val *reflect.Value) (nn int, err error) {
	return fmt.Fprintf(w, ":%d\n\r", val.Int())
}

func writeString(w io.Writer, val *reflect.Value) (nn int, err error) {
	return fmt.Fprintf(w, "+%s\n\r", val.String())
}

func writeError(w io.Writer, val *reflect.Value) (nn int, err error) {
	return fmt.Fprintf(w, "-%s\n\r", val.String())
}
