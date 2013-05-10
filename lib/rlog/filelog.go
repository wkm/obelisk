package rlog

import (
	"bufio"
	"fmt"
	"os"
)

type FileLog struct {
	File     *os.File
	Category string
	Prefix   string

	buff *bufio.Writer
}

func NewFileLog(fname, ctg, prefix string) (*FileLog, error) {
	log := FileLog{}
	file, err := os.OpenFile(fname, os.O_APPEND, 0770)
	if err != nil {
		return nil, err
	}

	log.File = file
	log.buff = bufio.NewWriter(log.File)
	log.Category = ctg
	log.Prefix = prefix

	return &log, nil
}

func (f FileLog) Printf(format string, obj ...interface{}) {
	lines, _ := fmt.Fprintf(f.buff, format, obj...)
	f.buff.Write([]byte("\n"))
	statByte.Add(uint(lines) + 1)
}

func (f FileLog) Sync() {
	f.buff.Flush()
}

func (f FileLog) Close() {
	f.Close()
}
