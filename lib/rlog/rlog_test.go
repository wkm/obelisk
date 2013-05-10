package rlog

import (
	"testing"
)

func TestLog(t *testing.T) {
	c := NewConfig()
	log := c.Logger("foo")

	prox, ok := log.(*LogProxy)
	if !ok {
		t.Errorf("expected log to be a LogProxy but was %#v", log)
	}

	if prox.Delegate == nil {
		t.Fatalf("expected delegate to exist")
	}

	delegate := *prox.Delegate
	_, ok = delegate.(StdoutLog)
	if !ok {
		t.Errorf("expected delegate to be a StdoutLog but was %#v", delegate)
	}
}
