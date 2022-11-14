package double

import (
	"fmt"
	"os"
	"strings"

	"github.com/cxcn/dtool/pkg/pinyin"
	"github.com/cxcn/dtool/pkg/table"
	"gopkg.in/ini.v1"
)

// 双拼映射表
type mapping struct {
	shengmu map[string]string
	yunmu   map[string]string
	yinjie  map[string]string
	rule    int
}

const (
	AABC = iota
	ABCC
)

func ToDoublePinyin(dict pinyin.Dict, path string, rule int) table.Table {
	m := newMapping(path, rule)
	ret := make(table.Table, 0, len(dict))
	for i := range dict {
		ret = append(ret, table.Entry{
			Word: dict[i].Word,
			Code: m.match(dict[i].Pinyin),
		})
	}
	return ret
}

func newMapping(path string, rule int) *mapping {
	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	m := &mapping{
		shengmu: cfg.Section("ShengMu").KeysHash(),
		yunmu:   cfg.Section("YunMu").KeysHash(),
		yinjie:  cfg.Section("YinJie").KeysHash(),
		rule:    rule,
	}
	toLower(m.shengmu)
	toLower(m.yunmu)
	toLower(m.yinjie)
	return m
}

func (m *mapping) match(pinyin []string) string {
	var ret string
	switch {
	case len(pinyin) < 2:
		ret = m.get(pinyin[0])
	case len(pinyin) == 2:
		ret = m.get(pinyin[0]) + m.get(pinyin[1])
	case len(pinyin) >= 4:
		ret = string([]byte{
			m.get(pinyin[0])[0], m.get(pinyin[1])[0], m.get(pinyin[2])[0], m.get(pinyin[len(pinyin)-1])[0],
		})
	default:
		switch m.rule {
		case AABC:
			ret = m.get(pinyin[0]) + string([]byte{m.get(pinyin[1])[0], m.get(pinyin[2])[0]})
		case ABCC:
			ret = string([]byte{m.get(pinyin[0])[0], m.get(pinyin[1])[0]}) + m.get(pinyin[2])
		}
	}
	return ret
}

func (m *mapping) get(yinjie string) string {
	if v, ok := m.yinjie[yinjie]; ok {
		return v
	} else if len(yinjie) < 2 {
		return v + "##"
	}

	var sm, ym string
	switch yinjie[0] {
	case 'a', 'o', 'e':
		sm = "#"
		ym = yinjie
	default:
		if yinjie[1] == 'h' {
			sm = yinjie[:2]
			ym = yinjie[2:]
		} else {
			sm = yinjie[:1]
			ym = yinjie[1:]
		}
	}
	if tmp, ok := m.shengmu[sm]; ok {
		sm = tmp
	}
	if tmp, ok := m.yunmu[ym]; ok {
		ym = tmp
	}
	yj := sm + ym
	m.yinjie[yinjie] = yj

	return yj
}

func toLower(m map[string]string) {
	for k, v := range m {
		m[k] = strings.ToLower(v)
	}
}
