package pinyin

import (
	"strings"

	"github.com/cxcn/dtool/pkg/encoder"
	"github.com/cxcn/dtool/pkg/table"
)

type MspyDat struct{}

func (MspyDat) Parse(path string) Dict {
	t := table.MsUDP{}.Parse(path)
	ret := make(Dict, 0, len(t))
	for i := range t {
		code := encoder.GetPinyin(t[i].Word)
		if len(code) == 0 {
			code = []string{t[i].Code}
		}
		ret = append(ret, Entry{
			Word:   t[i].Word,
			Pinyin: code,
			Freq:   1,
		})
	}
	return ret
}

func (MspyDat) Gen(dict Dict) []byte {
	return table.MsUDP{}.Gen(DictToTable(dict, ""))
}

func DictToTable(dict Dict, sep string) table.Table {
	t := make(table.Table, 0, len(dict))
	for i := range dict {
		t = append(t, table.Entry{
			Word: dict[i].Word,
			Code: strings.Join(dict[i].Pinyin, sep),
		})
	}
	return t
}
