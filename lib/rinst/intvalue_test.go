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

	b := make(MeasurementBuffer, 0)
	i.Measure("fig", &b)
	if len(b) != 1 {
		t.Error("intvalue bad measure length")
	}

	recv := b[0]
	if recv.Name != "fig" || recv.IntValue != 15 {
		t.Errorf("intvalue bad measure value %#v", recv)
	}
}
