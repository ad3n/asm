package qsort

import "unsafe"

func unsafeBytesTo64(b []byte) []uint64 {
	return unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(b))), len(b)/8)
}

func unsafeBytesTo128(b []byte) []uint128 {
	return unsafe.Slice((*uint128)(unsafe.Pointer(unsafe.SliceData(b))), len(b)/16)
}

func unsafeBytesTo192(b []byte) []uint192 {
	return unsafe.Slice((*uint192)(unsafe.Pointer(unsafe.SliceData(b))), len(b)/24)
}

func unsafeBytesTo256(b []byte) []uint256 {
	return unsafe.Slice((*uint256)(unsafe.Pointer(unsafe.SliceData(b))), len(b)/32)
}
