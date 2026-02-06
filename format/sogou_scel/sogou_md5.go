package sogou_scel

import (
	"encoding/binary"
	"io"
	"math/bits"
)

func magicNumbers() [4]uint32 {
	return [4]uint32{
		0x67452301,
		0xefcdab89,
		0x98badcfe,
		0x10325476,
	}
}

func checkSumStream(rd io.Reader) [4]uint32 {
	state := magicNumbers()
	block := make([]byte, 64)
	length := 0

	// 每次读取 64 字节
	for {
		n, err := rd.Read(block)
		length += n
		if n != 64 || err != nil {
			block = block[:n]
			break
		} else {
			blockGeneric(&state, [64]byte(block))
		}
	}

	// 填充后计算
	padding := padBlock(block, length)
	n := len(padding) / 64
	for i := 0; i < n; i++ {
		blockGeneric(&state, [64]byte(padding[64*i:]))
	}
	return state
}

func checkSum(x []byte) [4]uint32 {
	state := magicNumbers()

	// 64 字节分块数量，除去最后一块
	n := len(x)/64 - 1

	// 先计算前面部分
	for i := 0; i < n; i++ {
		blockGeneric(&state, [64]byte(x[64*i:]))
	}

	// 填充后计算
	padding := padBlock(x[n*64:], len(x))
	n = len(padding) / 64
	for i := 0; i < n; i++ {
		blockGeneric(&state, [64]byte(padding[64*i:]))
	}
	return state
}

// 使用时最好只传入最后一块的数据和总的长度
//
// 填充到 64 的倍数，length 为原始数据字节数
func padBlock(padding []byte, length int) []byte {
	// （以位为单位）
	length <<= 3

	// 计算需要填充的字节数
	paddingLength := 64 - (len(padding)+8)%64
	if paddingLength < 0 {
		paddingLength += 64
	}

	// 创建填充后的数据
	padded := make([]byte, len(padding)+paddingLength+8)
	copy(padded, padding)

	// 添加 0x80 字节
	padded[len(padding)] = 0x80

	// 添加 0x00 字节
	for i := len(padding) + 1; i < len(padded)-8; i++ {
		padded[i] = 0x00
	}

	// 添加原始数据长度（64 位，小端序）
	binary.LittleEndian.PutUint64(padded[len(padded)-8:], uint64(length))

	return padded
}

func rotateLeft(x uint32, n uint) uint32 {
	return bits.RotateLeft32(x, int(n))
}

func blockGeneric(state *[4]uint32, block [64]byte) {
	var x [16]uint32
	for i := 0; i < 16; i++ {
		x[i] = binary.LittleEndian.Uint32(block[i*4:])
	}

	a, b, c, d := state[0], state[1], state[2], state[3]

	// round 1
	a = (^b&d | c&b) + x[0] + 0xD76AA478 + a
	a = rotateLeft(a, 7) + b
	d = (^a&c | b&a) + x[1] + 0xE8C7B756 + d
	d = rotateLeft(d, 12) + a
	c = (^d&b | d&a) + x[2] + 0x242070DB + c
	c = rotateLeft(c, 17) + d
	b = (^c&a | d&c) + x[3] + 0xC1BDCEEE + b
	b = rotateLeft(b, 22) + c
	a = (^b&d | c&b) + x[4] + 0xF57C0FAF + a
	a = rotateLeft(a, 7) + b
	d = (^a&c | b&a) + x[5] + 0x4787C62A + d
	d = rotateLeft(d, 12) + a
	c = (^d&b | d&a) + x[6] + 0xA8304613 + c
	c = rotateLeft(c, 17) + d
	b = (^c&a | d&c) + x[7] + 0xFD469501 + b
	b = rotateLeft(b, 22) + c
	a = (^b&d | c&b) + x[8] + 0x698098D8 + a
	a = rotateLeft(a, 7) + b
	d = (^a&c | b&a) + x[9] + 0x8B44F7AF + d
	d = rotateLeft(d, 12) + a
	c = (^d&b | d&a) + x[10] + 0xFFFF5BB1 + c
	c = rotateLeft(c, 17) + d
	b = (^c&a | d&c) + x[11] + 0x895CD7BE + b
	b = rotateLeft(b, 22) + c
	a = (^b&d | c&b) + x[12] + 0x6B901122 + a
	a = rotateLeft(a, 7) + b
	d = (^a&c | b&a) + x[13] + 0xFD987193 + d
	d = rotateLeft(d, 12) + a
	c = (^d&b | d&a) + x[14] + 0xA679438E + c
	c = rotateLeft(c, 17) + d
	b = (^c&a | d&c) + x[15] + 0x49B40821 + b
	b = rotateLeft(b, 22) + c

	// round 2
	a = (^d&c | d&b) + x[1] + 0xF61E2562 + a
	a = rotateLeft(a, 5) + b
	d = (^c&b | c&a) + x[6] + 0xC040B340 + d
	d = rotateLeft(d, 9) + a
	c = (^b&a | d&b) + x[11] + 0x265E5A51 + c
	c = rotateLeft(c, 14) + d
	b = (^a&d | c&a) + x[0] + 0xE9B6C7AA + b
	b = rotateLeft(b, 20) + c
	a = (^d&c | d&b) + x[5] + 0xD62F105D + a
	a = rotateLeft(a, 5) + b
	d = (^c&b | c&a) + x[10] + 0x02441453 + d
	d = rotateLeft(d, 9) + a
	c = (^b&a | d&b) + x[15] + 0xD8A1E681 + c
	c = rotateLeft(c, 14) + d
	b = (^a&d | c&a) + x[4] + 0xE7D3FBC8 + b
	b = rotateLeft(b, 20) + c
	a = (^d&c | d&b) + x[9] + 0x21E1CDE6 + a
	a = rotateLeft(a, 5) + b
	d = (^c&b | c&a) + x[14] + 0xC33707D6 + d
	d = rotateLeft(d, 9) + a
	c = (^b&a | d&b) + x[3] + 0xF4D50D87 + c
	c = rotateLeft(c, 14) + d
	b = (^a&d | c&a) + x[8] + 0x455A14ED + b
	b = rotateLeft(b, 20) + c
	a = (^d&c | d&b) + x[13] + 0xA9E3E905 + a
	a = rotateLeft(a, 5) + b
	d = (^c&b | c&a) + x[2] + 0xFCEFA3F8 + d
	d = rotateLeft(d, 9) + a
	c = (^b&a | d&b) + x[7] + 0x676F02D9 + c
	c = rotateLeft(c, 14) + d
	b = (^a&d | c&a) + x[12] + 0x8D2A4C8A + b
	b = rotateLeft(b, 20) + c

	// round 3
	a = (d ^ c ^ b) + x[5] + 0xFFFA3942 + a
	a = rotateLeft(a, 4) + b
	d = (c ^ b ^ a) + x[8] + 0x8771F681 + d
	d = rotateLeft(d, 11) + a
	c = (d ^ b ^ a) + x[11] + 0x6D9D6122 + c
	c = rotateLeft(c, 16) + d
	b = (d ^ c ^ a) + x[14] + 0xFDE5380C + b
	magic := rotateLeft(b, 23) + c
	a = a + 0xA4BEEA44 + (d ^ c ^ magic) + x[1]
	b = rotateLeft(a, 4) + magic
	d = (c ^ magic ^ b) + x[4] + 0x4BDECFA9 + d
	d = rotateLeft(d, 11) + b
	c = (d ^ magic ^ b) + x[7] + 0xF6BB4B60 + c
	c = rotateLeft(c, 16) + d
	a = magic + 0xBEBFBC70 + (d ^ c ^ b) + x[10]
	a = rotateLeft(a, 23) + c
	b = (d ^ c ^ a) + x[13] + 0x289B7EC6 + b
	b = rotateLeft(b, 4) + a
	d = (c ^ a ^ b) + x[0] + 0xEAA127FA + d
	d = rotateLeft(d, 11) + b
	c = (d ^ a ^ b) + x[3] + 0xD4EF3085 + c
	c = rotateLeft(c, 16) + d
	a = a + 0x04881D05 + (d ^ c ^ b) + x[6]
	a = rotateLeft(a, 23) + c
	b = (d ^ c ^ a) + x[9] + 0xD9D4D039 + b
	b = rotateLeft(b, 4) + a
	d = (c ^ a ^ b) + x[12] + 0xE6DB99E5 + d
	d = rotateLeft(d, 11) + b
	c = (d ^ a ^ b) + x[15] + 0x1FA27CF8 + c
	c = rotateLeft(c, 16) + d
	a = (d ^ c ^ b) + x[2] + 0xC4AC5665 + a
	a = rotateLeft(a, 23) + c

	// round 4
	b = ((^d | a) ^ c) + x[0] + 0xF4292244 + b
	b = rotateLeft(b, 6) + a
	d = ((^c | b) ^ a) + x[7] + 0x432AFF97 + d
	d = rotateLeft(d, 10) + b
	c = ((^a | d) ^ b) + x[14] + 0xAB9423A7 + c
	c = rotateLeft(c, 15) + d
	a = ((^b | c) ^ d) + x[5] + 0xFC93A039 + a
	a = rotateLeft(a, 21) + c
	b = ((^d | a) ^ c) + x[12] + 0x655B59C3 + b
	b = rotateLeft(b, 6) + a
	d = ((^c | b) ^ a) + x[3] + 0x8F0CCC92 + d
	d = rotateLeft(d, 10) + b
	c = ((^a | d) ^ b) + x[10] + 0xFFEFF47D + c
	c = rotateLeft(c, 15) + d
	a = ((^b | c) ^ d) + x[1] + 0x85845DD1 + a
	a = rotateLeft(a, 21) + c
	b = ((^d | a) ^ c) + x[8] + 0x6FA87E4F + b
	b = rotateLeft(b, 6) + a
	d = ((^c | b) ^ a) + x[15] + 0xFE2CE6E0 + d
	d = rotateLeft(d, 10) + b
	c = ((^a | d) ^ b) + x[6] + 0xA3014314 + c
	c = rotateLeft(c, 15) + d
	a = ((^b | c) ^ d) + x[13] + 0x4E0811A1 + a
	a = rotateLeft(a, 21) + c
	b = ((^d | a) ^ c) + x[4] + 0xF7537E82 + b
	b = rotateLeft(b, 6) + a
	d = ((^c | b) ^ a) + x[11] + 0xBD3AF235 + d
	d = rotateLeft(d, 10) + b
	c = ((^a | d) ^ b) + x[2] + 0x2AD7D2BB + c
	c = rotateLeft(c, 15) + d
	a = ((^b | c) ^ d) + x[9] + 0xEB86D391 + a
	a = rotateLeft(a, 21) + c

	state[0] += b
	state[1] += a
	state[2] += c
	state[3] += d
}
