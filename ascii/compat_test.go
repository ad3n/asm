package ascii

import (
	"bytes"
	"strings"
	"testing"
)

func TestEqualFoldASCIICompatibility(t *testing.T) {
	for a := 0; a < 128; a++ {
		for b := 0; b < 128; b++ {
			for _, n := range []int{1, 7, 8, 15, 16, 31, 32, 65, 128} {
				left, right := bytes.Repeat([]byte{byte(a)}, n), bytes.Repeat([]byte{byte(b)}, n)
				if EqualFold(left, right) != bytes.EqualFold(left, right) || EqualFoldString(string(left), string(right)) != strings.EqualFold(string(left), string(right)) {
					t.Fatalf("ASCII fold mismatch a=%d b=%d length=%d", a, b, n)
				}
			}
		}
	}
}
