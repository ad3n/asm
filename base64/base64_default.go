//go:build purego || !(amd64 || arm64)
// +build purego !amd64,!arm64

package base64

import "encoding/base64"

type Encoding = base64.Encoding

func newEncoding(encoder string) *Encoding {
	return base64.NewEncoding(encoder)
}

func fromStandardEncoding(enc *base64.Encoding) *Encoding {
	return enc
}

func standardEncoding(enc *Encoding) *base64.Encoding {
	return enc
}
