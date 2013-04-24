package rlog

import (
	"bytes"
	"circuit/use/circuit"
	"fmt"
)

type RLog struct {
	Buffer bytes.Buffer
}

const ServiceName = "remote-log"

// a global object containing the main log
var Log = new(RLog)

func init() {
	circuit.RegisterValue(&RLog{})
}

// print a message into the remote log
func (r *RLog) Printf(format string, obj ...interface{}) {
	fmt.Fprintf(&r.Buffer, format, obj...)
}

// gets the current content of the log; truncating the contents
func (r *RLog) FlushLog() []byte {
	content := r.Buffer.Bytes()
	r.Buffer.Reset()
	return content
}

// circuit.Listen(ServiceName, Log)
