//go:build !purego

package mem

import "unsafe"

func Copy(dst, src []byte) int {
	n := min(len(dst), len(src))
	distance := uintptr(unsafe.Pointer(unsafe.SliceData(dst))) - uintptr(unsafe.Pointer(unsafe.SliceData(src)))
	if distance < uintptr(n) {
		return copy(dst, src)
	}

	return copyAsm(dst, src)
}
