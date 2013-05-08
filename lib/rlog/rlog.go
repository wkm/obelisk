/*
	logging facility for workers
*/
package rlog

import (
	"bytes"
	"fmt"
)

type Log interface {
	// write a message to the log
	Printf(format string, obj ...interface{})

	// ensure the log's buffer has been synchronized
	Sync()
}

type MemoryLog struct {
	Buffer bytes.Buffer
}

// print a message into the log, terminated by a new line
func (r *MemoryLog) Printf(format string, obj ...interface{}) {
	statPrint.Incr()
	lines, _ := fmt.Fprintf(&r.Buffer, format, obj...)
	r.Buffer.WriteRune('\n')
	statByte.Add(uint(lines) + 1)
}

func (r *MemoryLog) Sync() {
	// memory logs don't synchronize
}

// gets the current content of the log; truncating the contents
func (r *MemoryLog) FlushLog() []byte {
	statFlush.Incr()
	content := r.Buffer.Bytes()
	r.Buffer.Reset()
	return content
}

// create a logger for a named category
func Logger(ctg string) Log {
	return nil
}
