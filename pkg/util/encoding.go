package util

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
)

// 将 io流 转换为 utf-8
func NewReader(input io.Reader) io.Reader {
	brd := bufio.NewReader(input)
	buf, _ := brd.Peek(1024)
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf) // 检测编码格式
	if err != nil {
		return brd
	}
	if cs.Confidence != 100 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}
	// 删除 BOM 文件头
	boms := make(map[string][]byte)
	boms["UTF-16BE"] = []byte{0xfe, 0xff}
	boms["UTF-16LE"] = []byte{0xff, 0xfe}
	boms["UTF-8"] = []byte{0xef, 0xbb, 0xbf}
	if b, ok := boms[cs.Charset]; ok {
		if bytes.HasPrefix(buf, b) {
			brd.Read(b)
		}
	}
	rd, _ := charset.NewReaderLabel(cs.Charset, brd) // 转换字节流
	return rd
}

// 读取文件
func Read(path string) (io.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return f, err
	}
	return NewReader(f), nil
}

// []byte, encoding -> string
func Decode(b []byte, e string) (string, error) {
	decoder := new(encoding.Decoder)
	switch e {
	case "utf16":
		decoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	case "gbk":
		decoder = simplifiedchinese.GBK.NewDecoder()
	default:
		return "", errors.New("error encoding format")
	}
	ret, err := decoder.Bytes(b)
	return string(ret), err
}

// string, encoding -> []byte
func Encode(s string, e string) ([]byte, error) {
	encoder := new(encoding.Encoder)
	switch e {
	case "utf16":
		encoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	case "gbk":
		encoder = simplifiedchinese.GBK.NewEncoder()
	default:
		return []byte{}, errors.New("error encoding format")
	}
	ret, err := encoder.Bytes([]byte(s))
	return ret, err
}
