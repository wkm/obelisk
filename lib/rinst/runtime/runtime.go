package runtime

import (
	"fmt"
	"obelisk/lib/rinst"
	"runtime"
	"time"
)

var Stats = rinst.NewCollection()

func str(r uint64) string {
	return fmt.Sprintf("%d", r)
}

var statsGauge = rinst.GaugeValue{
	// MeasureFn
	func(n string, b rinst.MeasurementBuffer) {
		now := uint64(time.Now().Unix())
		r := runtime.MemStats{}
		runtime.ReadMemStats(&r)

		b <- rinst.Measurement{n + "alloc", now, str(r.Alloc)}
		b <- rinst.Measurement{n + "totalAlloc", now, str(r.TotalAlloc)}
		b <- rinst.Measurement{n + "sys", now, str(r.Sys)}
		b <- rinst.Measurement{n + "lookup", now, str(r.Lookups)}
		b <- rinst.Measurement{n + "malloc", now, str(r.Mallocs)}
		b <- rinst.Measurement{n + "free", now, str(r.Frees)}
	},

	// SchemaFn
	func(n string, b rinst.SchemaBuffer) {
		// general statistics
		b <- rinst.Schema{n + "alloc", rinst.TypeValue, "byte", "bytes allocated and still in use"}
		b <- rinst.Schema{n + "totalAlloc", rinst.TypeValue, "byte", "bytes allocated (even if freed)"}
		b <- rinst.Schema{n + "sys", rinst.TypeValue, "byte", "bytes obtained from the system"}
		b <- rinst.Schema{n + "lookup", rinst.TypeCounter, "lookup", "number of pointer lookups"}
		b <- rinst.Schema{n + "malloc", rinst.TypeCounter, "malloc", "number of mallocs"}
		b <- rinst.Schema{n + "free", rinst.TypeCounter, "free", "number of frees"}

		// heap statistics
		b <- rinst.Schema{n + "heap.alloc", rinst.TypeValue, "byte", "bytes allocated and still in use"}
		b <- rinst.Schema{n + "heap.sys", rinst.TypeValue, "byte", "bytes obtained from the system"}
		b <- rinst.Schema{n + "heap.idle", rinst.TypeValue, "byte", "bytes in idle spans"}
		b <- rinst.Schema{n + "heap.inuse", rinst.TypeValue, "byte", "bytes in non-idle spans"}
		b <- rinst.Schema{n + "heap.released", rinst.TypeValue, "byte", "bytes released to the OS"}
		b <- rinst.Schema{n + "heap.objects", rinst.TypeValue, "byte", "total number of allocated objects"}
	},
}
