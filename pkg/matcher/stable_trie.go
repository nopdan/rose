package matcher

import (
	"sort"
)

// 稳定 trie 树
type stableTrie struct {
	ch   map[rune]*stableTrie
	code string
	pos  int

	line uint32
}

func NewStableTrie() *stableTrie {
	t := new(stableTrie)
	t.ch = make(map[rune]*stableTrie, 1000)
	return t
}

var orderLine uint32 = 0

func (t *stableTrie) Insert(word, code string, pos int) {
	for _, v := range word {
		if t.ch == nil {
			t.ch = make(map[rune]*stableTrie)
			t.ch[v] = new(stableTrie)
		} else if t.ch[v] == nil {
			t.ch[v] = new(stableTrie)
		}
		t = t.ch[v]
	}
	// 同一个词取码表位置靠前的
	if t.code == "" {
		t.code = code
		t.pos = pos
		orderLine++
		t.line = orderLine
	}
}

func (t *stableTrie) Build() {
}

// 前缀树按码表序匹配
func (t *stableTrie) Match(text []rune) (int, string, int) {
	type res_tmp struct {
		wordLen int
		code    string
		pos     int
		line    uint32
	}
	res := make([]res_tmp, 0, 10)
	for p := 0; p < len(text); {
		t = t.ch[text[p]]
		p++
		if t == nil {
			break
		}
		if t.code != "" {
			res = append(res, res_tmp{p, t.code, t.pos, t.line})
		}
	}
	if len(res) == 0 {
		return 0, "", 1
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].line < res[j].line
	})
	return res[0].wordLen, res[0].code, res[0].pos
}
