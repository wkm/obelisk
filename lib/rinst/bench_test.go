package rinst

import (
	"testing"
)

// Benchmark atomic longs
func BenchmarkCounter(b *testing.B) {
	c := Counter{}
	for i := 0; i < b.N; i++ {
		c.Incr()
	}
}

// An int value is just two atomic things
func BenchmarkIntValue(b *testing.B) {
	v := IntValue{}
	for i := 0; i < b.N; i++ {
		v.Set(int64(i))
	}
}

// A floating value is mostly a mutex
func BenchmarkFloatValue(b *testing.B) {
	v := FloatValue{}
	for i := 0; i < b.N; i++ {
		v.Set(float64(i))
	}
}
