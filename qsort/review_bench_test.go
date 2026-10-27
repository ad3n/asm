package qsort

import (
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

func BenchmarkReview(b *testing.B) {
	for _, size := range []int{8, 16, 24, 32} {
		src := make([]byte, size*100)
		rand.New(rand.NewSource(1)).Read(src)
		for _, impl := range []struct {
			name string
			fn   func([]byte, int)
		}{{"asm", func(b []byte, n int) { Sort(b, n, nil) }}, {"stdlib", func(b []byte, n int) { sort.Sort(newGeneric(b, n, nil)) }}} {
			b.Run(strconv.Itoa(size)+"/"+impl.name, func(b *testing.B) {
				data := make([]byte, len(src))
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					copy(data, src)
					impl.fn(data, size)
				}
			})
		}
	}
}
