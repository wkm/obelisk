package agent

import (
	"fmt"
	"time"

	"github.com/cloudfoundry/gosigar"
	"github.com/wkm/obelisk/lib/rinst"
)

func fstr(f float64) string { return fmt.Sprintf("%f", f) }
func istr(i uint64) string  { return fmt.Sprintf("%d", i) }

var StatsGauge = rinst.GaugeValue{
	// measure function
	func(n string, r rinst.MeasurementReceiver) {
		uptime := sigar.Uptime{}

		load := sigar.LoadAverage{}
		mem := sigar.Mem{}
		swap := sigar.Swap{}

		now := uint64(time.Now().Unix())
		load.Get()
		b <- rinst.Measurement{n + "load.1", now, fstr(load.One)}
		b <- rinst.Measurement{n + "load.5", now, fstr(load.Five)}
		b <- rinst.Measurement{n + "load.15", now, fstr(load.Fifteen)}

		mem.Get()
		b <- rinst.Measurement{n + "mem", now, istr(mem.Total)}
		b <- rinst.Measurement{n + "mem.used", now, istr(mem.Used)}
		b <- rinst.Measurement{n + "mem.actUsed", now, istr(mem.ActualUsed)}

		swap.Get()
		b <- rinst.Measurement{n + "swap", now, istr(swap.Total)}
		b <- rinst.Measurement{n + "swap.used", now, istr(swap.Used)}

		uptime.Get()
		b <- rinst.Measurement{n + "uptime", now, fstr(uptime.Length)}
	},
	// schema function
	func(n string, b rinst.SchemaBuffer) {
		b <- rinst.Schema{n + "load.1", rinst.TypeFloatValue, "proc", "one minute load"}
		b <- rinst.Schema{n + "load.5", rinst.TypeFloatValue, "proc", "five minute load"}
		b <- rinst.Schema{n + "load.15", rinst.TypeFloatValue, "proc", "fifteen minute load"}

		b <- rinst.Schema{n + "mem", rinst.TypeAllocation, "byte", "system memory usage"}
		b <- rinst.Schema{n + "mem.used", rinst.TypeFloatValue, "byte", "memory used"}
		b <- rinst.Schema{n + "mem.actUsed", rinst.TypeFloatValue, "byte", "actual memory used"}

		b <- rinst.Schema{n + "swap", rinst.TypeAllocation, "byte", "swap memory usage"}
		b <- rinst.Schema{n + "swap.used", rinst.TypeFloatValue, "byte", "swap used"}

	},
}
