package qsort

import (
	"bytes"
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

func TestSortGeneric(t *testing.T) {
	for _, size := range []int{1, 3, 7, 15, 33, 255, 256, 257, 1024} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			data := make([]byte, size*100)
			rand.New(rand.NewSource(1)).Read(data)
			original := bytes.Clone(data)
			indices := make([]int, 100)
			for i := range indices {
				indices[i] = i
			}

			Sort(data, size, func(i, j int) {
				indices[i], indices[j] = indices[j], indices[i]
			})
			if !sort.IsSorted(newGeneric(data, size, nil)) {
				t.Fatal("output is not sorted")
			}

			for i, index := range indices {
				if !bytes.Equal(data[i*size:(i+1)*size], original[index*size:(index+1)*size]) {
					t.Fatal("callback does not match output")
				}
			}
		})
	}
}

func BenchmarkSortGeneric(b *testing.B) {
	for _, size := range []int{3, 15, 33, 257, 1024} {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			original := make([]byte, size*100)
			rand.New(rand.NewSource(1)).Read(original)
			data := make([]byte, len(original))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				copy(data, original)
				Sort(data, size, nil)
			}
		})
	}
}
