package log

import (
	"runtime"
	"testing"
)

func TestAlloc(t *testing.T) {
	pool := NewBufferPool()

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	runtime.GC()

	for i := 0; i < 1000; i++ {
		b := pool.Get()
		b.WriteString("Hello")
		pool.Put(b)
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)

	frees := m2.Frees - m1.Frees
	if frees > 1000 {
		t.Fatalf("expected less than 100 frees after GC, got %d", frees)
	}
}

func BenchmarkAlloc(b *testing.B) {
	pool := NewBufferPool()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b := pool.Get()
		b.WriteString("Hello")
		pool.Put(b)
	}
}
