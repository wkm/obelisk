package resp

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
)

// Listen implements a simple read-execute-write dispatch loop against an interface
func Listen(reciever interface{}, r io.Reader, w io.Writer) (err error) {
	d, err := newDispatch(reciever)
	if err != nil {
		return
	}

	br := bufio.NewReader(r)
	bw := bufio.NewWriter(w)

	for {
		line, _, err := br.ReadLine()
		if err != nil {
			return err
		}

		out, callerr := d.Call(string(line))
		if callerr != nil {
			writeError(bw, callerr.Error())
		} else {
			bw.WriteString(out)
		}

		bw.Flush()
	}
}

// Write a request into the writer
func WriteRequest(w io.Writer, cmds ...interface{}) (nn int, err error) {
	// Send header
	cnn, err := fmt.Fprintf(w, "*%d\r\n", len(cmds))
	nn += cnn
	if err != nil {
		return
	}

	for _, c := range cmds {
		r := reflect.TypeOf(c)
		cnn, err = write(w, r.Kind(), reflect.ValueOf(c))
		nn += cnn

		if err != nil {
			return
		}
	}

	return
}
