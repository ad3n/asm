//go:build (amd64 || arm64) && !purego
// +build amd64 arm64
// +build !purego

package base64

import (
	"encoding/base64"
	"slices"

	"github.com/segmentio/asm/internal/unsafebytes"
)

var zeroEncoding base64.Encoding

type Encoding struct {
	enc    func(dst []byte, src []byte, lut *int8) (int, int)
	enclut [encLutSize]int8

	dec    func(dst []byte, src []byte, lut *int8) (int, int)
	declut [decLutSize]int8

	base *base64.Encoding
}

func (enc Encoding) WithPadding(padding rune) *Encoding {
	enc.base = standardEncoding(&enc).WithPadding(padding)
	return &enc
}

func (enc Encoding) Strict() *Encoding {
	enc.base = standardEncoding(&enc).Strict()
	return &enc
}

func (enc *Encoding) Encode(dst, src []byte) {
	if len(src) == 0 {
		return
	}

	if len(src) >= minEncodeLen && enc.enc != nil && len(dst) >= enc.base.EncodedLen(len(src)) {
		d, n := enc.enc(dst, src, &enc.enclut[0])
		dst = dst[d:]
		src = src[n:]
	}

	standardEncoding(enc).Encode(dst, src)
}

func (enc *Encoding) EncodeToString(src []byte) string {
	buf := make([]byte, standardEncoding(enc).EncodedLen(len(src)))
	enc.Encode(buf, src)
	return unsafebytes.String(buf)
}

func (enc *Encoding) EncodedLen(n int) int {
	return standardEncoding(enc).EncodedLen(n)
}

func (enc *Encoding) Decode(dst, src []byte) (n int, err error) {
	if len(src) < minDecodeLen || enc.dec == nil || len(dst) < len(src)/4*3 {
		if len(src) == 0 {
			return 0, nil
		}

		return standardEncoding(enc).Decode(dst, src)
	}

	d, consumed := enc.dec(dst, src, &enc.declut[0])
	n, err = enc.base.Decode(dst[d:], src[consumed:])
	if corrupt, ok := err.(base64.CorruptInputError); ok && consumed != 0 {
		err = corrupt + base64.CorruptInputError(consumed)
	}

	return n + d, err
}

func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	src := unsafebytes.BytesOf(s)
	dst := make([]byte, standardEncoding(enc).DecodedLen(len(s)))
	n, err := enc.Decode(dst, src)
	return dst[:n], err
}

func (enc *Encoding) DecodedLen(n int) int {
	return standardEncoding(enc).DecodedLen(n)
}

func (enc *Encoding) AppendEncode(dst, src []byte) []byte {
	n := enc.EncodedLen(len(src))
	dst = slices.Grow(dst, n)
	enc.Encode(dst[len(dst):][:n], src)
	return dst[:len(dst)+n]
}

func (enc *Encoding) AppendDecode(dst, src []byte) ([]byte, error) {
	return standardEncoding(enc).AppendDecode(dst, src)
}

func fromStandardEncoding(enc *base64.Encoding) *Encoding {
	return &Encoding{base: enc}
}

func standardEncoding(enc *Encoding) *base64.Encoding {
	if enc.base == nil {
		return &zeroEncoding
	}

	return enc.base
}
