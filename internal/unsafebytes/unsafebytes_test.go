package unsafebytes

import "testing"

func TestConversions(t *testing.T) {
	for _, b := range [][]byte{nil, {}, make([]byte, 0, 8), []byte("hello")} {
		if String(b) != string(b) {
			t.Fatal("string conversion changed contents")
		}

		if len(b) > 0 && Pointer(b) != &b[0] {
			t.Fatal("pointer does not reference input")
		}
	}

	for _, s := range []string{"", "hello"} {
		b := BytesOf(s)
		if string(b) != s || len(b) != len(s) || cap(b) != len(s) {
			t.Fatal("byte view changed contents or dimensions")
		}
	}
}

func TestConversionsAllocs(t *testing.T) {
	b := []byte("hello")
	if allocs := testing.AllocsPerRun(100, func() {
		if String(BytesOf(String(b))) != "hello" || Pointer(b) != &b[0] {
			panic("invalid conversion")
		}
	}); allocs != 0 {
		t.Fatalf("got %v allocations, want zero", allocs)
	}
}
