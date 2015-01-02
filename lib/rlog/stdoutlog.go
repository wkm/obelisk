package rlog

import (
	"fmt"
	"os"
)

// StdoutLog is a simple implementation of a log which writes directly to standard out.
type StdoutLog struct {
	Category string
}

// Printf writes a formatted log line.
func (l StdoutLog) Printf(format string, obj ...interface{}) {
	statPrint.Incr()
	params := make([]interface{}, len(obj)+1)
	params[0] = l.Category
	copy(params[1:], obj)
	bytes, _ := fmt.Fprintf(os.Stderr, "[%s] "+format+"\n", params...)
	statByte.Add(uint(bytes))
}

// Sync ensures all log lines have been saved.
func (l StdoutLog) Sync() {
	os.Stderr.Sync()
}

// Close disables all writing to the log.
func (l StdoutLog) Close() {
	// stdout logs don't close
}
