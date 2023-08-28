package util

import (
	"encoding/binary"
	"io"

	"golang.org/x/exp/constraints"
)

// 读取指定数量字节
func ReadN[T constraints.Integer](r io.Reader, size T) []byte {
	tmp := make([]byte, size)
	r.Read(tmp)
	return tmp
}

func ReadIntN[T constraints.Integer](r io.Reader, size T) int {
	tmp := ReadN(r, size)
	return BytesToInt(tmp)
}

// 读取小端 uint16
func ReadUint16(r io.Reader) uint16 {
	tmp := ReadN(r, 2)
	return binary.LittleEndian.Uint16(tmp)
}

// 读取小端 uint32
func ReadUint32(r io.Reader) uint32 {
	tmp := ReadN(r, 4)
	return binary.LittleEndian.Uint32(tmp)
}
