package wubi

import (
	"io"

	"github.com/nopdan/rose/util"
	"golang.org/x/exp/constraints"
)

var (
	ReadUint16 = util.ReadUint16
	ReadUint32 = util.ReadUint32
)

func ReadN[T constraints.Integer](r io.Reader, size T) []byte {
	return util.ReadN(r, size)
}

func ReadIntN[T constraints.Integer](r io.Reader, size T) int {
	return util.ReadIntN(r, size)
}

func To2Bytes[T constraints.Integer](i T) []byte {
	return util.To2Bytes(i)
}

func To4Bytes[T constraints.Integer](i T) []byte {
	return util.To4Bytes(i)
}
