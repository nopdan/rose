package wubi

import (
	"bufio"
	"bytes"
)

const WORDS = 3

type Words struct {
	*Custom
}

func init() {
	FormatList = append(FormatList, NewWords())
}
func NewWords() *Words {
	f := newCustom("t|w", "UTF-8")
	f.Name = "纯词组"
	f.ID = "words"
	f.CanMarshal = true
	return &Words{Custom: f}
}

func (f *Words) GetKind() int {
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
