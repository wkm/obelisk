package rinst

import (
	"testing"
)

func TestGauge(t *testing.T) {
	c := NewCollection()
	c.AddInstrument("foo", &GaugeValue{
		func(n string, r MeasurementReceiver) {
			r.WriteInt(n+"a", 123, 456)
		},

		func(n string, r SchemaReceiver) {
			r.WriteSchema(n+"a", TypeCounter, "ops", "an ops")
		},
	})

	mb := make(MeasurementBuffer, 0)
	c.Measure("", &mb)

	if len(mb) != 1 {
		t.Errorf("layout snapshot bad length %d", len(mb))
	}

	recv := mb[0]
	if recv.Name != "foo.a" || recv.IntValue != 456 {
		t.Errorf("layout snapshot bad value %#v", recv)
	}

	sb := make(SchemaBuffer, 0)
	c.Schema("", &sb)

	if len(sb) != 1 {
		t.Errorf("layout schema bad length %d", len(sb))
	}

	recv2 := sb[0]
	if recv2.Name != "foo.a" || recv2.Description != "an ops" || recv2.Unit != "ops" {
		t.Errorf("layout schema bad value %#v", recv2)
	}
}
