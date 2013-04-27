// logging facility for workers, available remotely
package rlog

import (
	"bytes"
	"fmt"
)

type MemoryLog struct {
	Buffer bytes.Buffer
}

// print a message into the log, terminated by a new line
func (r *MemoryLog) Printf(format string, obj ...interface{}) {
	statPrints.Incr()
	lines, _ := fmt.Fprintf(&r.Buffer, format, obj...)
	r.Buffer.WriteRune('\n')
	statBytes.Add(uint(lines) + 1)
}

// gets the current content of the log; truncating the contents
func (r *MemoryLog) FlushLog() []byte {
	statFlushes.Incr()
	content := r.Buffer.Bytes()
	r.Buffer.Reset()
	return content
}
