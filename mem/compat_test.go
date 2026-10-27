package mem

import (
	"bytes"
	"testing"
)

func TestCopyOverlapCompatibility(t *testing.T) {
	for _, size := range []int{0, 1, 15, 16, 31, 32, 63, 64, 65, 127, 128, 129, 256, 1024, 4096} {
		for _, offset := range []int{-65, -17, -1, 0, 1, 17, 65} {
			got := make([]byte, size+256)
			for i := range got {
				got[i] = byte(i * 31)
			}

			want := bytes.Clone(got)
			gn := Copy(got[128+offset:128+offset+size], got[128:128+size])
			wn := copy(want[128+offset:128+offset+size], want[128:128+size])
			if gn != wn || !bytes.Equal(got, want) {
				t.Fatalf("copy mismatch size=%d offset=%d", size, offset)
			}
		}
	}
}

func TestContainsByteCompatibility(t *testing.T) {
	for n := 0; n < 260; n++ {
		data := make([]byte, n)
		for i := range data {
			data[i] = byte(i * 31)
		}

		for value := 0; value < 256; value++ {
			if ContainsByte(data, byte(value)) != bytes.Contains(data, []byte{byte(value)}) {
				t.Fatalf("ContainsByte mismatch length=%d byte=%d", n, value)
			}
		}
	}
}

func TestCopyStackAllocs(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		var dst, src [128]byte
		src[0] = 42
		if Copy(dst[:], src[:]) != len(src) || dst[0] != src[0] {
			panic("copy failed")
		}
	})
	if allocs != 0 {
		t.Fatalf("got %v allocations, want zero", allocs)
	}
}

func TestBlendStackAllocs(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		var dst, src [128]byte
		src[0] = 42
		if Blend(dst[:], src[:]) != len(src) || dst[0] != src[0] {
			panic("blend failed")
		}
	})
	if allocs != 0 {
		t.Fatalf("got %v allocations, want zero", allocs)
	}
}

func TestMaskStackAllocs(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		var dst, src [128]byte
		dst[0], src[0] = 255, 42
		if Mask(dst[:], src[:]) != len(src) || dst[0] != src[0] {
			panic("mask failed")
		}
	})
	if allocs != 0 {
		t.Fatalf("got %v allocations, want zero", allocs)
	}
}
