package rinst

import (
	"testing"
)

func TestFloatValue(t *testing.T) {
	v := &FloatValue{}
	v.Set(-12)
	v.Set(13)
	v.Set(14)

	if v.Get() != 14 {
		t.Error("floatvalue bad value")
	}

	if v.NumSets() != 3 {
		t.Error("floatvalue bad numchanges")
	}

	b := make(MeasurementBuffer, 10)
	v.Measure("fig", b)
	if len(b) != 2 {
		t.Error("floatvalue bad measure length")
	}

	expec := Measurement{"fig", 0, "14.000000"}
	actual := <-b
	if expec != actual {
		t.Errorf("floatvalue bad measure value %v", actual)
	}

	expec = Measurement{"fig.sets", 0, "3"}
	actual = <-b
	if expec != actual {
		t.Error("floatvalue bad measure numsets %v", actual)
	}
}
