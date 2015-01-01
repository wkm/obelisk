package rlog

import (
	"bufio"
	"fmt"
	"os"
)

// FileLog is a log which is backed by a file.
type FileLog struct {
	File     *os.File
	Category string
	Prefix   string

	buff *bufio.Writer
}

// NewFileLog creates a log instance which is file backed.
func NewFileLog(fname, ctg, prefix string) (*FileLog, error) {
	log := FileLog{}
	file, err := os.OpenFile(fname, os.O_APPEND, 0770)
	if err != nil {
		return nil, err
	}

	log.File = file
	log.Category = ctg
	log.Prefix = prefix

	log.buff = bufio.NewWriter(log.File)

	return &log, nil
}

// Printf writes a line into the log.
func (f FileLog) Printf(format string, obj ...interface{}) {
	lines, _ := fmt.Fprintf(f.buff, format, obj...)
	f.buff.Write([]byte("\n"))
	statByte.Add(uint(lines) + 1)
}

// Sync ensures all log statements have been persisted.
func (f FileLog) Sync() {
	f.buff.Flush()
}

// Close closes the log and the underlying file.
func (f FileLog) Close() {
	f.Close()
}
