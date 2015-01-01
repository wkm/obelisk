package runtime

import (
	"runtime"
	"time"

	"github.com/wkm/obelisk/lib/rinst"
)

// Stats provides measurements on the Go runtime.
var Stats = &statsGauge

var statsGauge = rinst.GaugeValue{
	MeasureFn: func(n string, r rinst.MeasurementReceiver) {
		now := time.Now().Unix()
		s := runtime.MemStats{}
		runtime.ReadMemStats(&s)

		r.WriteInt(n+"alloc", now, int64(s.Alloc))
		r.WriteInt(n+"totalAlloc", now, int64(s.TotalAlloc))
		r.WriteInt(n+"sys", now, int64(s.Sys))
		r.WriteInt(n+"lookup", now, int64(s.Lookups))
		r.WriteInt(n+"malloc", now, int64(s.Mallocs))
		r.WriteInt(n+"free", now, int64(s.Frees))

		r.WriteInt(n+"heap.alloc", now, int64(s.HeapAlloc))
		r.WriteInt(n+"heap.sys", now, int64(s.HeapSys))
		r.WriteInt(n+"heap.idle", now, int64(s.HeapIdle))
		r.WriteInt(n+"heap.inuse", now, int64(s.HeapInuse))
		r.WriteInt(n+"heap.released", now, int64(s.HeapReleased))
		r.WriteInt(n+"heap.objects", now, int64(s.HeapObjects))

		r.WriteInt(n+"stack.inuse", now, int64(s.StackInuse))
		r.WriteInt(n+"stack.sys", now, int64(s.StackSys))
		r.WriteInt(n+"mspan.inuse", now, int64(s.MSpanInuse))
		r.WriteInt(n+"mspan.sys", now, int64(s.MSpanSys))
		r.WriteInt(n+"mcache.inuse", now, int64(s.MCacheInuse))
		r.WriteInt(n+"mcache.sys", now, int64(s.MCacheSys))
		r.WriteInt(n+"buckhashsys", now, int64(s.BuckHashSys))

		r.WriteInt(n+"gc.num", now, int64(s.NumGC))
		r.WriteInt(n+"gc.date", now, int64(s.LastGC))

		r.WriteInt(n+"gc.pause", now, int64(s.PauseTotalNs))

		// Boolean Flags
		// r.WriteInt(n+"gc.enabled", now, int64(s.EnableGC))
		// r.WriteInt(n+"gc.debugenabled", now, bools.DebugGC)
	},

	SchemaFn: func(n string, r rinst.SchemaReceiver) {
		// general statistics
		r.WriteSchema(n+"alloc", rinst.TypeIntValue, "byte", "bytes allocated and still in use")
		r.WriteSchema(n+"totalAlloc", rinst.TypeIntValue, "byte", "bytes allocated (even if freed)")
		r.WriteSchema(n+"sys", rinst.TypeIntValue, "byte", "bytes obtained from the system")
		r.WriteSchema(n+"lookup", rinst.TypeCounter, "lookup", "number of pointer lookups")
		r.WriteSchema(n+"malloc", rinst.TypeCounter, "malloc", "number of mallocs")
		r.WriteSchema(n+"free", rinst.TypeCounter, "free", "number of frees")

		// heap statistics
		r.WriteSchema(n+"heap.alloc", rinst.TypeIntValue, "byte", "bytes allocated and still in use")
		r.WriteSchema(n+"heap.sys", rinst.TypeIntValue, "byte", "bytes obtained from the system")
		r.WriteSchema(n+"heap.idle", rinst.TypeIntValue, "byte", "bytes in idle spans")
		r.WriteSchema(n+"heap.inuse", rinst.TypeIntValue, "byte", "bytes in non-idle spans")
		r.WriteSchema(n+"heap.released", rinst.TypeIntValue, "byte", "bytes released to the OS")
		r.WriteSchema(n+"heap.objects", rinst.TypeIntValue, "byte", "total number of allocated objects")

		// fixed size allocation
		r.WriteSchema(n+"stack.inuse", rinst.TypeIntValue, "byte", "bootstrap stacks")
		r.WriteSchema(n+"stack.sys", rinst.TypeIntValue, "byte", "bootstrap stacks")
		r.WriteSchema(n+"mspan.inuse", rinst.TypeIntValue, "byte", "mspan structures")
		r.WriteSchema(n+"mspan.sys", rinst.TypeIntValue, "byte", "mspan structures")
		r.WriteSchema(n+"mcache.inuse", rinst.TypeIntValue, "byte", "mcache structures")
		r.WriteSchema(n+"mcache.sys", rinst.TypeIntValue, "byte", "mcache structures")
		r.WriteSchema(n+"buckhashsys", rinst.TypeIntValue, "byte", "profiling bucket hash table bytes obtained from system")

		r.WriteSchema(n+"gc.num", rinst.TypeCounter, "gc", "number of garbage collections")
		r.WriteSchema(n+"gc.date", rinst.TypeDateValue, "", "last run in absolute time")
		r.WriteSchema(n+"gc.pause", rinst.TypeCounter, "ns", "running run time in garbage collection")

		r.WriteSchema(n+"gc.enabled", rinst.TypeBoolValue, "", "garbage collection enabled")
		r.WriteSchema(n+"gc.debugenabled", rinst.TypeBoolValue, "", "debug garbage collection enabled")
	},
}
