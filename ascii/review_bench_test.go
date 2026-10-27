package ascii

import (
	"strconv"
	"strings"
	"testing"
)

var reviewResult bool

func BenchmarkReview(b *testing.B) {
	for _, n := range []int{32, 4096} {
		a, c := strings.Repeat("a", n), strings.Repeat("A", n)
		for _, impl := range []struct {
			name string
			fn   func(string, string) bool
		}{{"asm", EqualFoldString}, {"stdlib", strings.EqualFold}} {
			b.Run(strconv.Itoa(n)+"/"+impl.name, func(b *testing.B) {
				b.SetBytes(int64(n))
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					reviewResult = impl.fn(a, c)
				}
			})
		}
	}
}
