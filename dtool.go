package dtool

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	. "github.com/cxcn/dtool/utils"
)

// 单字码表 char: []code
type Dict map[rune][]string

// 初始化码表
func NewDict() Dict {
	return make(Dict)
}

// 生成编码
func (d Dict) Encode(s string, rules map[int][]int) map[string][]string {

	ret := make(map[string][]string)
	words := strings.Split(s, "\n")
	// 对一个词
	for _, w := range words {
		w = strings.TrimSpace(w)
		ret[w] = d.encodeWord(w, rules)
	}
	return ret
}

// 生成编码
func (d Dict) encodeWord(s string, rules map[int][]int) []string {
	word := []rune(s)
	if len(word) < 2 {
		log.Println("词条太短", s)
		return []string{}
	}
	rule, ok := rules[len(word)]
	if !ok { // 多字词的规则
		rule = rules[0]
	}

	var sliCode A
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
		var codes []string
		codes = d[char]
		// fmt.Println("char", string(char), "codes", codes)
		// 当前位置所有可能的字符
		tmp := make(map[byte]struct{})
		for _, code := range codes {
			var c byte
			if ci == 0 {
				c = code[len(code)-1]
			} else {
				if len(code) <= ci-1 {
					log.Println("索引超出", s, code, ci-1)
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
	ret := Product(sliCode)
	// fmt.Println(sliCode, ret)
	return ret
}

// 处理规则 AaAbBaBb,AaBaCaZa
func RuleHandle(rule string) map[int][]int {
	sliRule := strings.Split(rule, "\n")
	rules := make(map[int][]int)
	for _, r := range sliRule {
		tmp := strings.Split(r, "=")
		if len(tmp) != 2 {
			continue
		}
		id, err := strconv.Atoi(tmp[0])
		if err != nil {
			log.Fatal("规则解析错误")
		}
		// 把 AaAbBaBb 转为了 1 1 1 2 2 1 2 2
		rules[id] = func(s string) []int {
			s = strings.TrimSpace(s)
			// 65 97
			ret := make([]int, 0, 10)
			for i := range s {
				id := 0
				// Z and z
				if s[i] == 90 || s[i] == 122 {

				} else if s[i] >= 97 { // A-Z
					id = int(s[i]) - 96
				} else { // a-z
					id = int(s[i]) - 64
				}
				ret = append(ret, id)
			}
			return ret
		}(tmp[1])
	}
	return rules
}

// 读取码表，flag: 字在后
func (d Dict) Read(s string, flag bool) {
	f, err := os.Open(s)
	rd, err := Decode(f)
	if err != nil {
		log.Fatal("读取码表失败")
	}
	scan := bufio.NewScanner(rd)

out:
	for scan.Scan() {
		items := strings.Split(scan.Text(), "\t")
		if len(items) < 2 {
			log.Println("这一行没有以\\t分割", scan.Text())
			continue
		}
		var ws []rune
		var c string
		if flag { // 编码在前
			c, ws = items[0], []rune(items[1])
		} else { // 编码在后
			ws, c = []rune(items[0]), items[1]
		}
		if len(ws) != 1 {
			log.Println("单字长度不对,跳过", ws)
			continue
		}
		w := ws[0] // 是一个字
		// 去重
		if len(d[w]) != 0 {
			for i, v := range d[w] {
				if strings.HasPrefix(v, c) { // 原来的编码
					continue out
				} else if strings.HasPrefix(c, v) { // 当前的编码
					d[w][i] = c
					continue out
				}
			}
			d[w] = append(d[w], c)
		} else {
			d[w] = []string{c}
		}
	}
}
