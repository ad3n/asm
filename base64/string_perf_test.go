package base64

import (
	"bytes"
	stdbase64 "encoding/base64"
	"strconv"
	"testing"
)

var encodedStringSink string

func TestEncodeToStringOwnership(t *testing.T) {
	for _, n := range []int{0, 1, 2, 3, 32, 1024} {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			src := bytes.Repeat([]byte{0xa5}, n)
			got := StdEncoding.EncodeToString(src)
			want := stdbase64.StdEncoding.EncodeToString(src)
			clear(src)
			if got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
}

func BenchmarkEncodeToString(b *testing.B) {
	for _, n := range []int{32, 1024, 65536} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			src := bytes.Repeat([]byte{0xa5}, n)
			b.ReportAllocs()
			b.SetBytes(int64(n))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				encodedStringSink = StdEncoding.EncodeToString(src)
			}
		})
	}
}

func TestBufferOperationsAllocs(t *testing.T) {
	src := bytes.Repeat([]byte{0xa5}, 1024)
	encoded := make([]byte, StdEncoding.EncodedLen(len(src)))
	decoded := make([]byte, StdEncoding.DecodedLen(len(encoded)))
	allocs := testing.AllocsPerRun(100, func() {
		StdEncoding.Encode(encoded, src)
		n, err := StdEncoding.Decode(decoded, encoded)
		if err != nil || !bytes.Equal(decoded[:n], src) {
			panic("round trip failed")
		}
	})
	if allocs != 0 {
		t.Fatalf("got %v allocations, want zero", allocs)
	}
}
