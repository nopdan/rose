package wubi

import (
	"bufio"
	"bytes"
)

const WORDS = 3

type Words struct {
	Template
}

func init() {
	FormatList = append(FormatList, NewWords())
}
func NewWords() *Words {
	f := new(Words)
	f.Name = "纯词组"
	f.ID = "words"
	f.CanMarshal = true
	return f
}

func (f *Words) GetKind() int {
	return WORDS
}
func (f *Words) Unmarshal(r *bytes.Reader) []*Entry {
	return nil
}

// 特殊
func (f *Words) UnmarshalStr(r *bytes.Reader) []string {
	di := make([]string, 0, r.Size()>>8)
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		di = append(di, scan.Text())
	}
	return di
}

func (f *Words) MarshalStr(di []string) []byte {
	var buf bytes.Buffer
	for _, v := range di {
		buf.WriteString(v)
		buf.WriteString("\r\n")
	}
	return buf.Bytes()
}
