package rlog

import (
	"fmt"
	"os"
)

type StdoutLog struct {
	Category string
}

func (l StdoutLog) Printf(format string, obj ...interface{}) {
	statPrint.Incr()
	bytes, _ := fmt.Fprintf(os.Stderr, "[%s] ", l.Category)
	statByte.Add(uint(bytes))

	bytes, _ = fmt.Fprintf(os.Stderr, format, obj...)
	os.Stderr.WriteString("\n")
	statByte.Add(uint(bytes) + 1)
}

func (l StdoutLog) Sync() {
	os.Stderr.Sync()
}

func (l StdoutLog) Close() {
	// stdout logs don't close
}
