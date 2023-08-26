package wubi

import (
	"bufio"
	"bytes"
)

const WORDS = 3

type Words struct {
	*Custom
}

func NewWords() *Words {
	return &Words{Custom: NewCustom("t|w", "UTF-8")}
}

func (f *Words) GetType() int {
	return WORDS
}

func (f *Words) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)

	scan := bufio.NewScanner(r)
	for scan.Scan() {
		di = append(di, &Entry{scan.Text(), "", 0})
	}
	return di
}

// 特殊
func (f *Words) MarshalStr(di []string) []byte {
	var buf bytes.Buffer
	for _, v := range di {
		buf.WriteString(v)
		buf.WriteString("\r\n")
	}
	return buf.Bytes()
}
