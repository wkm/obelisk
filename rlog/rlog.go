// logging facility for workers, available remotely
package rlog

import (
	"bytes"
	"fmt"
)

type Log struct {
	Buffer bytes.Buffer
}

// print a message into the log
func (r *Log) Printf(format string, obj ...interface{}) {
	fmt.Fprintf(&r.Buffer, format, obj...)
}

// gets the current content of the log; truncating the contents
func (r *Log) FlushLog() []byte {
	content := r.Buffer.Bytes()
	r.Buffer.Reset()
	return content
}
