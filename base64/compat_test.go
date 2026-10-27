package base64

import (
	"bytes"
	stdbase64 "encoding/base64"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeCompatibility(t *testing.T) {
	for _, encoding := range encodings {
		for _, strict := range []bool{false, true} {
			candidate, control := encoding.candidate, encoding.control
			if strict {
				candidate, control = candidate.Strict(), control.Strict()
			}

			encoded := control.EncodeToString(bytes.Repeat([]byte{0xfb, 0xff, 0xef}, 64))
			for pos := 0; pos < len(encoded); pos++ {
				for _, replacement := range []byte{'!', '=', '+', '/', '-', '_', ',', '\r', '\n', 0xff} {
					src := []byte(encoded)
					src[pos] = replacement
					got := make([]byte, candidate.DecodedLen(len(src)))
					want := make([]byte, control.DecodedLen(len(src)))
					gn, ge := candidate.Decode(got, src)
					wn, we := control.Decode(want, src)
					if gn != wn || !reflect.DeepEqual(ge, we) || !bytes.Equal(got[:gn], want[:wn]) {
						t.Fatalf("%s strict=%v pos=%d replacement=%x: got (%d,%v), want (%d,%v)", encoding.name, strict, pos, replacement, gn, ge, wn, we)
					}
				}
			}

			for _, tail := range []string{"", "A", "AA", "AB", "AAA", "AAB", "AA==", "AB==", "AAA=", "AAB=", "\r\n"} {
				src := strings.Repeat("AAAA", 64) + tail
				got, ge := candidate.DecodeString(src)
				want, we := control.DecodeString(src)
				if !bytes.Equal(got, want) || !reflect.DeepEqual(ge, we) {
					t.Fatalf("%s strict=%v tail=%q: got %v, want %v", encoding.name, strict, tail, ge, we)
				}
			}
		}
	}
}

func TestEncodingMethodSet(t *testing.T) {
	got, want := reflect.TypeOf(StdEncoding), reflect.TypeOf(stdbase64.StdEncoding)
	if got.NumMethod() != want.NumMethod() {
		t.Fatalf("got %d methods, want %d", got.NumMethod(), want.NumMethod())
	}

	for i := 0; i < want.NumMethod(); i++ {
		method := want.Method(i)
		actual, ok := got.MethodByName(method.Name)
		if !ok || actual.Type.String() != method.Type.String() {
			t.Fatalf("incompatible method %s", method.Name)
		}
	}

	var err error = CorruptInputError(7)
	if _, ok := err.(stdbase64.CorruptInputError); !ok {
		t.Fatal("error type is not a stdlib alias")
	}
}

func TestExtendedCompatibility(t *testing.T) {
	alphabets := []string{encodeStd, encodeURL, encodeIMAP, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/"}
	high := make([]byte, 64)
	for i := range high {
		high[i] = byte(128 + i)
	}

	alphabets = append(alphabets, string(high))
	for _, alphabet := range alphabets {
		for _, padding := range []rune{StdPadding, NoPadding, '.', 255} {
			for _, strict := range []bool{false, true} {
				candidate := NewEncoding(alphabet).WithPadding(padding)
				control := stdbase64.NewEncoding(alphabet).WithPadding(padding)
				if strict {
					candidate, control = candidate.Strict(), control.Strict()
				}

				for n := 0; n < 260; n++ {
					src := bytes.Repeat([]byte{0xfb}, n)
					got, want := candidate.EncodeToString(src), control.EncodeToString(src)
					if got != want {
						t.Fatalf("encoding mismatch: padding=%d strict=%v n=%d", padding, strict, n)
					}

					for _, prefix := range [][]byte{nil, []byte("prefix")} {
						gd, wd := bytes.Clone(prefix), bytes.Clone(prefix)
						if !bytes.Equal(candidate.AppendEncode(gd, src), control.AppendEncode(wd, src)) {
							t.Fatal("AppendEncode mismatch")
						}

						gd, ge := candidate.AppendDecode(bytes.Clone(prefix), []byte(got))
						wd, we := control.AppendDecode(bytes.Clone(prefix), []byte(want))
						if !bytes.Equal(gd, wd) || !reflect.DeepEqual(ge, we) {
							t.Fatalf("AppendDecode mismatch: %v != %v", ge, we)
						}
					}

					decoded, ge := candidate.DecodeString(got)
					expected, we := control.DecodeString(want)
					if !bytes.Equal(decoded, expected) || !reflect.DeepEqual(ge, we) {
						t.Fatalf("DecodeString mismatch: %v != %v", ge, we)
					}
				}
			}
		}
	}
}

func TestAppendReuse(t *testing.T) {
	src := bytes.Repeat([]byte{0xa5}, 1024)
	encoded := make([]byte, 0, StdEncoding.EncodedLen(len(src)))
	encoded = StdEncoding.AppendEncode(encoded, src)
	decoded := make([]byte, 0, len(src))
	allocs := testing.AllocsPerRun(100, func() {
		encoded = StdEncoding.AppendEncode(encoded[:0], src)
		var err error
		decoded, err = StdEncoding.AppendDecode(decoded[:0], encoded)
		if err != nil || !bytes.Equal(decoded, src) {
			panic("round trip failed")
		}
	})
	if allocs != 0 {
		t.Fatalf("got %v allocations, want zero", allocs)
	}

	for _, src := range []string{"Zg==", "Zh==", "AAAA!", strings.Repeat("AAAA", 64) + "!", "Zg==\r\n"} {
		for _, enc := range encodings {
			got, ge := enc.candidate.Strict().AppendDecode([]byte("prefix"), []byte(src))
			want, we := enc.control.Strict().AppendDecode([]byte("prefix"), []byte(src))
			if !bytes.Equal(got, want) || !reflect.DeepEqual(ge, we) {
				t.Fatalf("invalid append input %q: got %v, want %v", src, ge, we)
			}
		}
	}
}

func TestStreamingCompatibility(t *testing.T) {
	src := bytes.Repeat([]byte("hello world"), 1000)
	for _, enc := range encodings {
		var got bytes.Buffer
		writer := NewEncoder(enc.candidate, &got)
		for start := 0; start < len(src); start += 7 {
			end := min(start+7, len(src))
			if _, err := writer.Write(src[start:end]); err != nil {
				t.Fatal(err)
			}
		}

		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}

		if got.String() != enc.control.EncodeToString(src) {
			t.Fatal("stream encoding mismatch")
		}

		decoded, err := io.ReadAll(NewDecoder(enc.candidate, bytes.NewReader(got.Bytes())))
		if err != nil || !bytes.Equal(decoded, src) {
			t.Fatalf("stream decoding mismatch: %v", err)
		}
	}
}

func TestZeroEncodingCompatibility(t *testing.T) {
	var candidate Encoding
	var control stdbase64.Encoding
	for _, src := range [][]byte{nil, []byte("hello"), bytes.Repeat([]byte("hello"), 100)} {
		got, want := candidate.EncodeToString(src), control.EncodeToString(src)
		if got != want {
			t.Fatal("zero encoding mismatch")
		}

		gd, ge := candidate.DecodeString(got)
		wd, we := control.DecodeString(want)
		if !bytes.Equal(gd, wd) || !reflect.DeepEqual(ge, we) {
			t.Fatal("zero decoding mismatch")
		}
	}
}

func FuzzDecodeCompatibility(f *testing.F) {
	for _, src := range []string{"", "Zg==", "Zh==", "\r\n", strings.Repeat("AAAA", 64) + "!", strings.Repeat("/+_-", 32)} {
		f.Add([]byte(src), uint8(0), false)
	}

	f.Fuzz(func(t *testing.T, src []byte, mode uint8, strict bool) {
		enc := encodings[int(mode)%len(encodings)]
		candidate, control := enc.candidate, enc.control
		if strict {
			candidate, control = candidate.Strict(), control.Strict()
		}

		got := make([]byte, candidate.DecodedLen(len(src)))
		want := make([]byte, control.DecodedLen(len(src)))
		gn, ge := candidate.Decode(got, src)
		wn, we := control.Decode(want, src)
		if gn != wn || !bytes.Equal(got[:gn], want[:wn]) || !reflect.DeepEqual(ge, we) {
			t.Fatalf("%s strict=%v: got (%d,%v), want (%d,%v)", enc.name, strict, gn, ge, wn, we)
		}
	})
}

func TestShortDestinationCompatibility(t *testing.T) {
	panics := func(fn func()) (result bool) {
		defer func() {
			result = recover() != nil
		}()

		fn()
		return false
	}
	for _, enc := range encodings {
		src := bytes.Repeat([]byte{0xa5}, 128)
		encoded := []byte(enc.control.EncodeToString(src))
		for _, encode := range []bool{false, true} {
			input, size := encoded, len(src)
			if encode {
				input, size = src, len(encoded)
			}

			for n := 0; n <= size; n++ {
				got, want := bytes.Repeat([]byte{0xcc}, size+32), bytes.Repeat([]byte{0xcc}, size+32)
				gp := panics(func() {
					if encode {
						enc.candidate.Encode(got[:n:n], input)
						return
					}

					enc.candidate.Decode(got[:n:n], input)
				})
				wp := panics(func() {
					if encode {
						enc.control.Encode(want[:n:n], input)
						return
					}

					enc.control.Decode(want[:n:n], input)
				})
				if gp != wp || !bytes.Equal(got, want) {
					t.Fatalf("%s encode=%v destination=%d: panic=%v want=%v or modified bytes differ", enc.name, encode, n, gp, wp)
				}
			}
		}
	}
}

func TestConstructorPanicsCompatibility(t *testing.T) {
	panicValue := func(fn func()) (value any) {
		defer func() {
			value = recover()
		}()

		fn()
		return nil
	}
	for _, alphabet := range []string{"", encodeStd[:63], encodeStd + "!", strings.Repeat("A", 64), "\n" + encodeStd[1:], "\r" + encodeStd[1:]} {
		got := panicValue(func() { NewEncoding(alphabet) })
		want := panicValue(func() { stdbase64.NewEncoding(alphabet) })
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("NewEncoding panic = %v, want %v", got, want)
		}
	}

	for _, padding := range []rune{-2, '\r', '\n', 'A', 256} {
		got := panicValue(func() { StdEncoding.WithPadding(padding) })
		want := panicValue(func() { stdbase64.StdEncoding.WithPadding(padding) })
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("WithPadding panic = %v, want %v", got, want)
		}
	}
}
