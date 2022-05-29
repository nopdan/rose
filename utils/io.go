package util

import (
	"bytes"
)

// 字节（小端）转为整数
func BytesToInt(b []byte) int {
	var ret int
	for i, v := range b {
		ret |= int(v) << (i << 3)
	}
	return ret
}

// 读 length 个字节，转为 int(倒着的)
func ReadInt(rd *bytes.Reader, length int) int {
	tmp := make([]byte, length)
	rd.Read(tmp)
	return BytesToInt(tmp)
}
