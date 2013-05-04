package rinst

import (
	"testing"
)

func TestGauge(t *testing.T) {
	c := NewCollection()
	c.AddInstrument("foo", &GaugeValue{
		func(n string, b MeasurementBuffer) {
			b <- Measurement{n + "a", 123, "thing"}
		},

		func(n string, b SchemaBuffer) {
			b <- Schema{n + "a", TypeCounter, "ops", "an ops"}
		},
	})

	mb := make(MeasurementBuffer, 10)
	c.Measure("", mb)

	if len(mb) != 1 {
		t.Errorf("layout snapshot bad length %d", len(mb))
	}

	recv := <-mb
	if recv.Name != "foo.a" || recv.Value != "thing" {
		t.Error("layout snapshot bad value %#v", recv)
	}

	sb := make(SchemaBuffer, 10)
	c.Schema("", sb)

	if len(sb) != 1 {
		t.Errorf("layout schema bad length %d", len(sb))
	}

	recv2 := <-sb
	if recv2.Name != "foo.a" || recv2.Description != "an ops" || recv2.Unit != "ops" {
		t.Errorf("layout schema bad value %#v", recv2)
	}
}
