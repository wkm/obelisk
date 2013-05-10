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
	params := make([]interface{}, len(obj)+1)
	params[0] = l.Category
	copy(params[1:], obj)
	bytes, _ := fmt.Fprintf(os.Stderr, "[%s] "+format+"\n", params...)
	statByte.Add(uint(bytes))
}

func (l StdoutLog) Sync() {
	os.Stderr.Sync()
}

func (l StdoutLog) Close() {
	// stdout logs don't close
}
