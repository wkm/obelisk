package rinst

import (
	"testing"
)

func TestFloatValue(t *testing.T) {
	t.Parallel()

	v := &FloatValue{}
	v.Set(-12)
	v.Set(13)
	v.Set(14)

	if v.Get() != 14 {
		t.Error("floatvalue bad value")
	}

	b := make(MeasurementBuffer, 0)
	v.Measure("fig", &b)
	if len(b) != 1 {
		t.Error("floatvalue bad measure length")
	}

	actual := b[0]
	if actual.Name != "fig" || actual.FloatValue != 14.0 {
		t.Errorf("floatvalue bad measure value %v", actual)
	}
}
