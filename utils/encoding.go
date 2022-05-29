package util

import (
	"bufio"
	"errors"
	"io"

	"github.com/gogs/chardet"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
)

// 将 io流 转换为 utf-8
func DecodeIO(f io.Reader) (io.Reader, error) {

	brd := bufio.NewReader(f)
	buf, _ := brd.Peek(1024)
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf) // 检测编码格式
	if err != nil {
		return brd, err
	}
	if cs.Confidence != 100 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}
	rd, err := charset.NewReaderLabel(cs.Charset, brd) // 转换字节流
	return rd, err
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
