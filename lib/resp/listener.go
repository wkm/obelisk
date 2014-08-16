package resp

import (
	"bufio"
	"io"
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
			// ... hmm what to do?
		}

		bw.WriteString(out)
		bw.Flush()
	}
}
