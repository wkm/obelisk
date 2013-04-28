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

	actual := <-b
	if actual.Name != "fig" || actual.Value != "14.000000" {
		t.Errorf("floatvalue bad measure value %v", actual)
	}

	actual = <-b
	if actual.Name != "fig.sets" || actual.Value != "3" {
		t.Error("floatvalue bad measure numsets %v", actual)
	}
}
