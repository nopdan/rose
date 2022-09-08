package checker

import (
	"log"
	"strings"

	"github.com/cxcn/dtool/pkg/util"
)

// 生成编码
func (c *Checker) Encode(s string) map[string][]string {
	ret := make(map[string][]string)
	words := strings.Split(s, "\n")
	// 对一个词
	for _, w := range words {
		w = strings.TrimSpace(w)
		ret[w] = c.EncodeWord(w)
	}
	return ret
}

// 生成单个词的编码
func (c *Checker) EncodeWord(s string) []string {
	word := []rune(s)
	if len(word) < 2 {
		log.Println("词条太短", s)
		return []string{}
	}
	rule, ok := c.Rule[len(word)]
	if !ok { // 多字词的规则
		rule = c.Rule[0]
	}

	var sliCode [][]byte
	for i := 0; i < len(rule)/2; i++ {
		wi := rule[2*i]   // word index
		ci := rule[2*i+1] // code index
		// 每个字的所有编码
		var char rune
		if wi == 0 {
			char = word[len(word)-1]
		} else {
			char = word[wi-1]
		}
		codes := c.Dict[char]
		// fmt.Println("char", string(char), "codes", codes)
		// 当前位置所有可能的字符
		tmp := make(map[byte]struct{})
		for _, code := range codes {
			var c byte
			if ci == 0 {
				c = code[len(code)-1]
			} else {
				if len(code) <= ci-1 {
					// log.Println("索引超出", s, code, ci-1)
					break
				}
				c = code[ci-1]
			}
			tmp[c] = struct{}{}
		}
		sliTmp := make([]byte, 0, len(tmp))
		for k := range tmp {
			sliTmp = append(sliTmp, k)
		}
		sliCode = append(sliCode, sliTmp)
	}
	tmp := util.Product(sliCode)
	ret := make([]string, len(tmp))
	for i := range tmp {
		ret[i] = string(tmp[i])
	}
	// fmt.Println(sliCode, ret)
	return ret
}
