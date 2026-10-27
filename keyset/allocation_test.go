package keyset

import "testing"

var lookupResult int

func BenchmarkLookupStackKey(b *testing.B) {
	set := New([][]byte{[]byte("key")})
	if set == nil {
		b.Skip("optimized lookup unavailable")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := [3]byte{'k', 'e', 'y'}
		lookupResult = Lookup(set, key[:])
	}
}

func TestLookupStackKeyAllocs(t *testing.T) {
	set := New([][]byte{[]byte("key")})
	if set == nil {
		t.Skip("optimized lookup unavailable")
	}

	allocs := testing.AllocsPerRun(100, func() {
		key := [3]byte{'k', 'e', 'y'}
		if Lookup(set, key[:]) != 0 {
			panic("key not found")
		}
	})
	if allocs != 0 {
		t.Fatalf("got %v allocations, want zero", allocs)
	}
}
