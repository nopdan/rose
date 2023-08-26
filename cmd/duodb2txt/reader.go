package main

import (
	"io"

	"github.com/nopdan/rose/util"
	"golang.org/x/exp/constraints"
)

var ReadUint32 = util.ReadUint32

func ReadN[T constraints.Integer](r io.Reader, size T) []byte {
	return util.ReadN(r, size)
}

func ReadIntN[T constraints.Integer](r io.Reader, size T) int {
	return util.ReadIntN(r, size)
}
