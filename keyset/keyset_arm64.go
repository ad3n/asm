//go:build !purego
// +build !purego

package keyset

//go:noescape
func Lookup(keyset []byte, key []byte) int
