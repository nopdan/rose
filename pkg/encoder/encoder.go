package encoder

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/nopdan/rose/pkg/util"
)

type Encoder interface {
	Encode(string) string // 编码一个词，可能有多个编码
}

func New(schema string, data []byte, isAABC bool) Encoder {
	if schema == "phrase" {
		return newPhrase()
	}
	if w := newWubi(schema, isAABC); w != nil {
		return w
	}

	var r io.Reader
	if data == nil {
		var err error
		r, err = util.Read(schema)
		if err != nil {
			fmt.Printf("读取文件失败：%s\n", schema)
			panic(err)
		}
	} else {
		rd := bytes.NewReader(data)
		r = util.NewReader(rd)
	}

	w := &Wubi{
		Char:   make(map[rune]string),
		IsAABC: isAABC,
	}
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		fields := strings.Split(scan.Text(), "\t")
		if len(fields) < 2 {
			continue
		}
		word := []rune(fields[0])
		// 跳过词组
		if len(word) != 1 {
			continue
		}
		char := word[0]
		w.Char[char] = fields[1]
	}
	return w
}

type Phrase struct {
	enc *Pinyin
}

func newPhrase() *Phrase {
	enc := NewPinyin()
	return &Phrase{enc: enc}
}

func (p Phrase) Encode(word string) string {
	return strings.Join(p.enc.Encode(word), "")
}
