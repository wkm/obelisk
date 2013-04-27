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

	expec := Measurement{"fig", "15"}
	if expec != <-b {
		t.Error("intvalue bad measure value")
	}

	expec = Measurement{"fig.sets", "3"}
	if expec != <-b {
		t.Error("intvalue bad measure numsets")
	}
}
