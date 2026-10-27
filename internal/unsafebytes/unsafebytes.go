package unsafebytes

import "unsafe"

func Pointer(b []byte) *byte {
	return unsafe.SliceData(b)
}

func String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func BytesOf(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
