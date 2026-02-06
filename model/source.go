package model

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/gogs/chardet"
	"github.com/nopdan/rose/util"
)

// ReadSeekCloser 统一可读可寻址可关闭
// 兼容 *os.File 与内存 Reader
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

// Source 统一输入来源
// 支持本地文件与内存数据
type Source interface {
	Name() string
	Open() (ReadSeekCloser, error)
	Bytes() ([]byte, error)
}

// FileSource 本地文件输入
// Reader 返回 *os.File（更高性能，支持 Seek）
type FileSource struct {
	Path string
}

func NewFileSource(path string) *FileSource {
	return &FileSource{Path: path}
}

func (s *FileSource) Name() string {
	return filepath.Base(s.Path)
}

func (s *FileSource) Open() (ReadSeekCloser, error) {
	return os.Open(s.Path)
}

func (s *FileSource) Bytes() ([]byte, error) {
	return os.ReadFile(s.Path)
}

// BytesSource 内存输入（Web上传等）
type BytesSource struct {
	name string
	data []byte
}

func NewBytesSource(name string, data []byte) *BytesSource {
	return &BytesSource{name: name, data: data}
}

func (s *BytesSource) Name() string {
	if s.name == "" {
		return "memory"
	}
	return s.name
}

func (s *BytesSource) Open() (ReadSeekCloser, error) {
	return &readSeekNopCloser{Reader: bytes.NewReader(s.data)}, nil
}

func (s *BytesSource) Bytes() ([]byte, error) {
	return s.data, nil
}

// OpenTextReader 打开文本输入并返回 UTF-8 reader 与编码名称
// 使用 chardet 检测编码，低置信度时默认 GB18030
func OpenTextReader(src Source) (io.Reader, string, func() error, error) {
	r, err := src.Open()
	if err != nil {
		return nil, "", nil, err
	}

	br := bufio.NewReader(r)
	buf, _ := br.Peek(1024)
	detector := chardet.NewTextDetector()
	cs, err := detector.DetectBest(buf)
	if err != nil {
		return br, "", r.Close, nil
	}
	if cs.Confidence < 95 && cs.Charset != "UTF-8" {
		cs.Charset = "GB18030"
	}

	enc := util.LookupEnc(cs.Charset)
	if enc == nil {
		return stripBOM(br), cs.Charset, r.Close, nil
	}

	return stripBOM(enc.NewDecoder().Reader(br)), cs.Charset, r.Close, nil
}

// stripBOM 去除常见 BOM（UTF-8/UTF-16/UTF-32）
func stripBOM(r io.Reader) io.Reader {
	br := bufio.NewReader(r)
	if b, err := br.Peek(4); err == nil {
		switch {
		case len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF:
			_, _ = br.Discard(3)
		case len(b) >= 4 && b[0] == 0x00 && b[1] == 0x00 && b[2] == 0xFE && b[3] == 0xFF:
			_, _ = br.Discard(4)
		case len(b) >= 4 && b[0] == 0xFF && b[1] == 0xFE && b[2] == 0x00 && b[3] == 0x00:
			_, _ = br.Discard(4)
		case len(b) >= 2 && b[0] == 0xFE && b[1] == 0xFF:
			_, _ = br.Discard(2)
		case len(b) >= 2 && b[0] == 0xFF && b[1] == 0xFE:
			_, _ = br.Discard(2)
		}
	}
	return br
}

type readSeekNopCloser struct {
	*bytes.Reader
}

func (r *readSeekNopCloser) Close() error {
	return nil
}

// DecodeBinary 使用统一 Source 构建 Reader 并执行解析
func DecodeBinary(src Source, parse func(*Reader) ([]*Entry, error)) ([]*Entry, error) {
	r, err := NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}
	return parse(r)
}
