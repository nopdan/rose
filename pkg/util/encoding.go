package util

import (
	"bufio"
	"bytes"
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
	enc := getEncoding(e)
	b, err := enc.NewDecoder().Bytes(b)
	return string(b), err
}

// string, encoding -> []byte
func Encode(b []byte, e string) ([]byte, error) {
	enc := getEncoding(e)
	return enc.NewEncoder().Bytes(b)
}

func getEncoding(enc string) encoding.Encoding {
	var encoding encoding.Encoding
	switch enc {
	case "UTF-16LE":
		encoding = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	case "UTF-16BE":
		encoding = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case "GBK":
		encoding = simplifiedchinese.GBK
	case "GB18030":
		encoding = simplifiedchinese.GB18030
	default:
		encoding = unicode.UTF8
	}
	// fmt.Println(enc, encoding)
	return encoding
}
