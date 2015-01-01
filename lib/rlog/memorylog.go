package rlog

import (
	"bytes"
	"fmt"
)

// MemoryLog is an in-memory log container.
type MemoryLog struct {
	Buffer *bytes.Buffer
}

// NewMemoryLog allocates a new in-memory log.
func NewMemoryLog() *MemoryLog {
	m := MemoryLog{}
	m.Buffer = new(bytes.Buffer)
	return &m
}

// Printf a message into the log, terminated by a new line
func (r MemoryLog) Printf(format string, obj ...interface{}) {
	statPrint.Incr()
	bytes, _ := fmt.Fprintf(r.Buffer, format, obj...)
	r.Buffer.WriteRune('\n')
	statByte.Add(uint(bytes) + 1)
}

// Sync is a NOP with memory logs.
func (r MemoryLog) Sync() {
	// memory logs don't synchronize
}

// Close is a NOP with memory logs.
func (r MemoryLog) Close() {
	// memory logs don't close
}

// FlushLog gets the current content of the log, truncating its contents.
func (r MemoryLog) FlushLog() []byte {
	statFlush.Incr()
	content := r.Buffer.Bytes()
	r.Buffer.Reset()
	return content
}
