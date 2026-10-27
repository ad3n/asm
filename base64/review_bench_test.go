package base64

import (
	"bytes"
	std "encoding/base64"
	"strconv"
	"testing"
)

type reviewEncoding interface {
	Encode([]byte, []byte)
	Decode([]byte, []byte) (int, error)
	EncodeToString([]byte) string
	EncodedLen(int) int
	DecodedLen(int) int
}

var reviewString string

func BenchmarkReview(b *testing.B) {
	for _, alphabet := range []struct{ name, alphabet string }{{"std", encodeStd}, {"url", encodeURL}, {"imap", encodeIMAP}} {
		for _, size := range []int{32, 4096} {
			src := bytes.Repeat([]byte{0xfb}, size)
			for _, impl := range []struct {
				name string
				enc  reviewEncoding
			}{{"asm", NewEncoding(alphabet.alphabet)}, {"stdlib", std.NewEncoding(alphabet.alphabet)}} {
				enc := impl.enc
				prefix := alphabet.name + "/" + strconv.Itoa(size) + "/" + impl.name
				b.Run(prefix+"/Encode", func(b *testing.B) {
					dst := make([]byte, enc.EncodedLen(len(src)))
					b.ReportAllocs()
					b.SetBytes(int64(len(src)))
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						enc.Encode(dst, src)
					}
				})
				encoded := []byte(enc.EncodeToString(src))
				b.Run(prefix+"/Decode", func(b *testing.B) {
					dst := make([]byte, enc.DecodedLen(len(encoded)))
					b.ReportAllocs()
					b.SetBytes(int64(len(src)))
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						enc.Decode(dst, encoded)
					}
				})
				b.Run(prefix+"/String", func(b *testing.B) {
					b.ReportAllocs()
					b.SetBytes(int64(len(src)))
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						reviewString = enc.EncodeToString(src)
					}
				})
			}
		}
	}
}
