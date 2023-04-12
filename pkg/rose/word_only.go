package rose

import (
	"bufio"
	"bytes"

	util "github.com/flowerime/goutil"
	"github.com/flowerime/rose/pkg/zhuyin"
)

type WordOnly struct{ Dict }

func NewWordOnly() *WordOnly {
	d := new(WordOnly)
	d.IsPinyin = true
	d.IsBinary = false
	d.Name = "纯汉字.txt"
	return d
}

func (d *WordOnly) GetDict() *Dict {
	return &d.Dict
}

func (d *WordOnly) Parse() {
	pyt := make(PyTable, 0, 0xff)
	scan := bufio.NewScanner(d.rd)
	for scan.Scan() {
		pyt = append(pyt, &PinyinEntry{scan.Text(), zhuyin.Get(scan.Text()), 1})
	}
	d.pyt = pyt
}

func (d *WordOnly) GenFrom(src *Dict) []byte {
	var buf bytes.Buffer
	for _, v := range src.pyt {
		buf.WriteString(v.Word)
		buf.WriteString(util.LineBreak)
	}
	return buf.Bytes()
}
