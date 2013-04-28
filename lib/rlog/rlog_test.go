package rlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	log := &MemoryLog{}
	log.Printf("%s", "hello")
	log.Printf("%d", 12)
	content := log.FlushLog()

	if string(content) != "hello\n12\n" {
		t.Errorf("invalid log content (%s)", string(content))
	}

	if statPrint.Value() != 2 {
		t.Errorf("invalid prints count %d", statPrint.Value())
	}
}
