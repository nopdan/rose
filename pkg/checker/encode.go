package checker

import (
	"strings"

	"github.com/imetool/dtool/pkg/util"
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
	for _, word := range words {
		word = strings.Split(word, "\t")[0]
		word = strings.TrimSpace(word)
		ret = append(ret, Entry{word, c.EncodeWord(word)})
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
		for _, char := range word {
			codes := c.Dict[char]
			for _, idx := range c.RuleZ {
				tmp = append(tmp, currCodes(codes, int(idx)))
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
	return ret
}

// 定长规则
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
		ret = append(ret, currCodes(codes, int(ci)))
		// fmt.Println("char", string(char), "codes", codes)
	}
	return ret
}

// 当前位置所有可能的字符
func currCodes(codes []string, idx int) []byte {
	tmp := make([]byte, 0, 1)
	for _, code := range codes {
		var b byte
		if idx >= len(code) {
			b = code[len(code)-1]
		} else {
			b = code[idx]
		}
		tmp = append(tmp, b)
	}
	tmp = util.RmRepeat(tmp)
	return tmp
}
