package base64

import (
	"encoding/base64"
	"io"
)

const (
	StdPadding  rune = base64.StdPadding
	NoPadding   rune = base64.NoPadding
	encodeStd        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	encodeURL        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	encodeIMAP       = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+,"
	letterRange      = int8('Z' - 'A' + 1)
)

var (
	StdEncoding    = NewEncoding(encodeStd)
	URLEncoding    = NewEncoding(encodeURL)
	RawStdEncoding = StdEncoding.WithPadding(NoPadding)
	RawURLEncoding = URLEncoding.WithPadding(NoPadding)
)

type CorruptInputError = base64.CorruptInputError

func NewEncoding(encoder string) *Encoding {
	switch encoder {
	case encodeStd, encodeURL, encodeIMAP:
		return newEncoding(encoder)
	default:
		return fromStandardEncoding(base64.NewEncoding(encoder))
	}
}

func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser {
	if enc == nil {
		return base64.NewEncoder(nil, w)
	}

	return base64.NewEncoder(standardEncoding(enc), w)
}

func NewDecoder(enc *Encoding, r io.Reader) io.Reader {
	if enc == nil {
		return base64.NewDecoder(nil, r)
	}

	return base64.NewDecoder(standardEncoding(enc), r)
}
