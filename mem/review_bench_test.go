package mem

import (
	"bytes"
	"strconv"
	"testing"
)

var reviewResult int

func BenchmarkReview(b *testing.B) {
	for _, n := range []int{32, 4096} {
		src, dst := bytes.Repeat([]byte{42}, n), make([]byte, n)
		for _, impl := range []struct {
			name string
			fn   func([]byte, []byte) int
		}{{"asm", Copy}, {"stdlib", copyGeneric}} {
			b.Run(strconv.Itoa(n)+"/"+impl.name, func(b *testing.B) {
				b.SetBytes(int64(n))
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					reviewResult = impl.fn(dst, src)
				}
			})
		}
	}
}
