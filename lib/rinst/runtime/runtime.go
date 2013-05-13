package runtime

import (
	"fmt"
	"obelisk/lib/rinst"
	"runtime"
	"time"
)

var Stats = &statsGauge

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

		b <- rinst.Measurement{n + "heap.alloc", now, str(r.HeapAlloc)}
		b <- rinst.Measurement{n + "heap.sys", now, str(r.HeapSys)}
		b <- rinst.Measurement{n + "heap.idle", now, str(r.HeapIdle)}
		b <- rinst.Measurement{n + "heap.inuse", now, str(r.HeapInuse)}
		b <- rinst.Measurement{n + "heap.released", now, str(r.HeapReleased)}
		b <- rinst.Measurement{n + "heap.objects", now, str(r.HeapObjects)}

		b <- rinst.Measurement{n + "stack.inuse", now, str(r.StackInuse)}
		b <- rinst.Measurement{n + "stack.sys", now, str(r.StackSys)}
		b <- rinst.Measurement{n + "mspan.inuse", now, str(r.MSpanInuse)}
		b <- rinst.Measurement{n + "mspan.sys", now, str(r.MSpanSys)}
		b <- rinst.Measurement{n + "mcache.inuse", now, str(r.MCacheInuse)}
		b <- rinst.Measurement{n + "mcache.sys", now, str(r.MCacheSys)}
		b <- rinst.Measurement{n + "buckhashsys", now, str(r.BuckHashSys)}

		b <- rinst.Measurement{n + "gc.num", now, str(uint64(r.NumGC))}
	},

	// SchemaFn
	func(n string, b rinst.SchemaBuffer) {
		// general statistics
		b <- rinst.Schema{n + "alloc", rinst.TypeIntValue, "byte", "bytes allocated and still in use"}
		b <- rinst.Schema{n + "totalAlloc", rinst.TypeIntValue, "byte", "bytes allocated (even if freed)"}
		b <- rinst.Schema{n + "sys", rinst.TypeIntValue, "byte", "bytes obtained from the system"}
		b <- rinst.Schema{n + "lookup", rinst.TypeCounter, "lookup", "number of pointer lookups"}
		b <- rinst.Schema{n + "malloc", rinst.TypeCounter, "malloc", "number of mallocs"}
		b <- rinst.Schema{n + "free", rinst.TypeCounter, "free", "number of frees"}

		// heap statistics
		b <- rinst.Schema{n + "heap.alloc", rinst.TypeIntValue, "byte", "bytes allocated and still in use"}
		b <- rinst.Schema{n + "heap.sys", rinst.TypeIntValue, "byte", "bytes obtained from the system"}
		b <- rinst.Schema{n + "heap.idle", rinst.TypeIntValue, "byte", "bytes in idle spans"}
		b <- rinst.Schema{n + "heap.inuse", rinst.TypeIntValue, "byte", "bytes in non-idle spans"}
		b <- rinst.Schema{n + "heap.released", rinst.TypeIntValue, "byte", "bytes released to the OS"}
		b <- rinst.Schema{n + "heap.objects", rinst.TypeIntValue, "byte", "total number of allocated objects"}

		// fixed size allocation
		b <- rinst.Schema{n + "stack.inuse", rinst.TypeIntValue, "byte", "bootstrap stacks"}
		b <- rinst.Schema{n + "stack.sys", rinst.TypeIntValue, "byte", "bootstrap stacks"}
		b <- rinst.Schema{n + "mspan.inuse", rinst.TypeIntValue, "byte", "mspan structures"}
		b <- rinst.Schema{n + "mspan.sys", rinst.TypeIntValue, "byte", "mspan structures"}
		b <- rinst.Schema{n + "mcache.inuse", rinst.TypeIntValue, "byte", "mcache structures"}
		b <- rinst.Schema{n + "mcache.sys", rinst.TypeIntValue, "byte", "mcache structures"}
		b <- rinst.Schema{n + "buckhashsys", rinst.TypeIntValue, "byte", "profiling bucket hash table bytes obtained from system"}

		b <- rinst.Schema{n + "gc.num", rinst.TypeCounter, "gc", "number of garbage collections"}
	},
}
