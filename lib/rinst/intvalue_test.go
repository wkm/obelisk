package rinst

import (
	"testing"
)

func TestIntValue(t *testing.T) {
	i := &IntValue{}
	i.Set(12)
	i.Set(13)
	i.Set(15)

	if i.Get() != 15 {
		t.Error("intvalue bad value")
	}

	if i.NumSets() != 3 {
		t.Error("intvalue bad numchanges")
	}

	b := make(MeasurementBuffer, 10)
	i.Measure("fig", b)
	if len(b) != 2 {
		t.Error("intvalue bad measure length")
	}

	recv := <-b
	if recv.Name != "fig" || recv.Value != "15" {
		t.Error("intvalue bad measure value %#v", recv)
	}

	recv = <-b
	if recv.Name != "fig.sets" || recv.Value != "3" {
		t.Error("intvalue bad measure numsets %#v", recv)
	}
}
