package rlog

import (
	"bytes"
	"fmt"
)

type MemoryLog struct {
	Buffer *bytes.Buffer
}

func NewMemoryLog() *MemoryLog {
	m := MemoryLog{}
	m.Buffer = new(bytes.Buffer)
	return &m
}

// print a message into the log, terminated by a new line
func (r MemoryLog) Printf(format string, obj ...interface{}) {
	statPrint.Incr()
	bytes, _ := fmt.Fprintf(r.Buffer, format, obj...)
	r.Buffer.WriteRune('\n')
	statByte.Add(uint(bytes) + 1)
}

func (r MemoryLog) Sync() {
	// memory logs don't synchronize
}

func (r MemoryLog) Close() {
	// memory logs don't close
}

// gets the current content of the log; truncating the contents
func (r MemoryLog) FlushLog() []byte {
	statFlush.Incr()
	content := r.Buffer.Bytes()
	r.Buffer.Reset()
	return content
}
