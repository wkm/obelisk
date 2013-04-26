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

	if i.NumChanges() != 3 {
		t.Error("intvalue bad numchanges")
	}

	if len(i.Measure()) != 2 {
		t.Error("intvalue bad measure length")
	}

	if i.Measure()[0] != "15" {
		t.Error("intvalue bad measure value")
	}

	if i.Measure()[1] != "3" {
		t.Error("intvalue bad measure numchanges")
	}
}
