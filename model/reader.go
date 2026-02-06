package model

import (
	"bytes"
	"os"

	"github.com/nopdan/rose/util"
)

type Reader struct {
	*bytes.Reader
	buffer []byte
}

func NewReader(r *bytes.Reader) *Reader {
	return &Reader{
		Reader: r,
		buffer: make([]byte, 8),
	}
}

func Open(name string) (*Reader, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return NewReader(bytes.NewReader(data)), nil
}

// NewReaderFromSource 从统一输入来源构建 Reader
func NewReaderFromSource(src Source) (*Reader, error) {
	data, err := src.Bytes()
	if err != nil {
		return nil, err
	}
	return NewReader(bytes.NewReader(data)), nil
}

// Read 读取指定长度的字节
func (r *Reader) ReadN(size int) []byte {
	buf := make([]byte, size)
	_, err := r.Read(buf)
	if err != nil {
		return nil
	}
	return buf
}

// ReadIntN 读取指定字节数的整数 1<= size <= 8
func (r *Reader) ReadIntN(size int) int {
	if size < 1 || size > len(r.buffer) {
		return 0
	}
	_, err := r.Read(r.buffer[:size])
	if err != nil {
		return 0
	}
	var result int
	for i := range size {
		result |= int(r.buffer[i]) << (i * 8)
	}
	return result
}

// ReadUint32 读取 32 位无符号整数
func (r *Reader) ReadUint32() uint32 {
	r.Read(r.buffer[:4])
	return uint32(r.buffer[0]) | uint32(r.buffer[1])<<8 | uint32(r.buffer[2])<<16 | uint32(r.buffer[3])<<24
}

// ReadUint16 读取 16 位无符号整数
func (r *Reader) ReadUint16() uint16 {
	r.Read(r.buffer[:2])
	return uint16(r.buffer[0]) | uint16(r.buffer[1])<<8
}

func (r *Reader) readWithBuffer(size int) []byte {
	if len(r.buffer) < size {
		r.buffer = make([]byte, size)
	}
	_, err := r.Read(r.buffer[:size])
	if err != nil {
		return nil
	}
	return r.buffer[:size]
}

// ReadString 读取指定长度的 UTF-8 字符串
func (r *Reader) ReadString(size int) string {
	b := r.readWithBuffer(size)
	if len(b) < size {
		return ""
	}
	return string(b)
}

// ReadStringEnc 读取指定长度的字符串，使用指定编码解码
func (r *Reader) ReadStringEnc(size int, enc *util.Encoding) string {
	b := r.readWithBuffer(size)
	if len(b) < size {
		return ""
	}
	return enc.Decode(b)
}
