package agent

import (
	"time"

	"github.com/cloudfoundry/gosigar"

	"github.com/wkm/obelisk/lib/rinst"
)

var statsGauge = rinst.GaugeValue{
	MeasureFn: func(n string, r rinst.MeasurementReceiver) {
		uptime := sigar.Uptime{}

		load := sigar.LoadAverage{}
		mem := sigar.Mem{}
		swap := sigar.Swap{}

		now := time.Now().Unix()
		load.Get()

		r.WriteFloat(n+"load.1", now, load.One)
		r.WriteFloat(n+"load.5", now, load.Five)
		r.WriteFloat(n+"load.15", now, load.Fifteen)

		mem.Get()
		r.WriteInt(n+"mem", now, int64(mem.Total))
		r.WriteInt(n+"mem.used", now, int64(mem.Used))
		r.WriteInt(n+"mem.actUsed", now, int64(mem.ActualUsed))

		swap.Get()
		r.WriteInt(n+"swap", now, int64(swap.Total))
		r.WriteInt(n+"swap.used", now, int64(swap.Used))

		uptime.Get()
		r.WriteFloat(n+"uptime", now, uptime.Length)
	},
	SchemaFn: func(n string, r rinst.SchemaReceiver) {
		r.WriteSchema(n+"load.1", rinst.TypeFloatValue, "proc", "one minute load")
		r.WriteSchema(n+"load.5", rinst.TypeFloatValue, "proc", "five minute load")
		r.WriteSchema(n+"load.15", rinst.TypeFloatValue, "proc", "fifteen minute load")

		r.WriteSchema(n+"mem", rinst.TypeAllocation, "byte", "system memory usage")
		r.WriteSchema(n+"mem.used", rinst.TypeFloatValue, "byte", "memory used")
		r.WriteSchema(n+"mem.actUsed", rinst.TypeFloatValue, "byte", "actual memory used")

		r.WriteSchema(n+"swap", rinst.TypeAllocation, "byte", "swap memory usage")
		r.WriteSchema(n+"swap.used", rinst.TypeFloatValue, "byte", "swap used")
	},
}
