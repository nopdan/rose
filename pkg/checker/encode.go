package checker

import (
	"strings"

	"github.com/cxcn/dtool/pkg/util"
)

type Entry struct {
	Word  string
	Codes []string
}

// 生成编码
func (c *Checker) Encode(s string) []Entry {
	ret := make([]Entry, 0, len(s)>>2)
	words := strings.Split(s, "\n")
	// 对一个词
	for _, w := range words {
		w = strings.Split(w, "\t")[0]
		w = strings.TrimSpace(w)
		ret = append(ret, Entry{w, c.EncodeWord(w)})
	}
	return ret
}

// 生成单个词的编码
func (c *Checker) EncodeWord(s string) []string {
	word := []rune(s)
	if len(word) < 2 {
		return []string{}
	}
	rule, ok := c.Rule[len(word)]
	if !ok { // 多字词的规则
		rule = c.Rule[0]
	}

	tmp := make([][]byte, 0, 1)
	if len(c.RuleZ) != 0 {
		for _, s := range word {
			codes := c.Dict[s]
			for _, idx := range c.RuleZ {
				ctmp := make([]byte, 0, 1)
				for _, code := range codes {
					var c byte
					if int(idx) >= len(code) {
						c = code[len(code)-1]
					} else {
						c = code[idx]
					}
					ctmp = append(ctmp, c)
				}
				ctmp = util.RmRepeat(ctmp)
				tmp = append(tmp, ctmp)
			}
		}
	} else {
		tmp = c.getCodes(rule, word)
	}

	tmp = util.Product(tmp)
	ret := make([]string, len(tmp))
	for i := range tmp {
		ret[i] = string(tmp[i])
	}
	// fmt.Println(sliCode, ret)
	return ret
}

func (c *Checker) getCodes(rule []byte, word []rune) [][]byte {
	ret := make([][]byte, 0, 1)
	for i := 0; i < len(rule)/2; i++ {
		wi := rule[2*i]   // word index
		ci := rule[2*i+1] // code index

		// 每个字的所有编码
		var char rune
		if int(wi) >= len(word) {
			char = word[len(word)-1]
		} else {
			char = word[wi]
		}
		codes := c.Dict[char]
		// fmt.Println("char", string(char), "codes", codes)

		// 当前位置所有可能的字符
		tmp := make([]byte, 0, 1)
		for _, code := range codes {
			var c byte
			if int(ci) >= len(code) {
				c = code[len(code)-1]
			} else {
				c = code[ci]
			}
			tmp = append(tmp, c)
		}
		tmp = util.RmRepeat(tmp)
		ret = append(ret, tmp)
	}
	return ret
}
