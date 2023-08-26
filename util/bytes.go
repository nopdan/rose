package util

import (
	"encoding/binary"

	"golang.org/x/exp/constraints"
)

// 字节（小端）转为整数
func BytesToInt(b []byte) int {
	var ret int
	for i, v := range b {
		ret |= int(v) << (i * 8)
	}
	return ret
}

// 转换为长度为 4 的字节切片
func To4Bytes[T constraints.Integer](i T) []byte {
	ret := make([]byte, 4)
	binary.LittleEndian.PutUint32(ret, uint32(i))
	return ret
}

// 转换为长度为 2 的字节切片
func To2Bytes[T constraints.Integer](i T) []byte {
	ret := make([]byte, 2)
	binary.LittleEndian.PutUint16(ret, uint16(i))
	return ret
}
