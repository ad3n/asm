//go:build ignore
// +build ignore

package main

import (
	. "github.com/mmcloughlin/avo/build"
	"github.com/ad3n/asm/build/internal/x86"
)

func init() {
	ConstraintExpr("!purego")
}

func main() {
	x86.GenerateCopy("copyAsm", "copies src to dst, returning the number of bytes written.", nil)
}
