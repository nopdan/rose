package checker

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/imetool/dtool/pkg/table"
	"github.com/imetool/goutil/util"
)

type Checker struct {
	Dict  map[rune][]string
	Rule  map[int][]byte
	RuleZ []byte
}

func NewChecker(path, rule string) *Checker {
	c := new(Checker)
	c.Dict = readDict(path)
	c.Rule, c.RuleZ = parseRule(rule)
	return c
}

func (c *Checker) Check(table table.Table) []byte {
	var buf bytes.Buffer
	buf.WriteString("词\t编码\t可能的编码\n")
out:
	// 遍历整个码表
	for i := range table {
		word, code := table[i].Word, table[i].Code
		// 排除单字
		if len([]rune(word)) == 1 {
			continue
		}
		// 根据规则推测可能的编码
		codes := c.EncodeWord(word)
		// 对每个编码
		for _, c := range codes {
			if strings.HasPrefix(c, code) {
				continue out
			}
		}
		if len(codes) != 0 {
			buf.WriteString(word + "\t" + code + "\t")
			buf.WriteString(strings.Join(codes, " "))
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
}

// 读取码表
func readDict(path string) map[rune][]string {
	rd, err := util.Read(path)
	if err != nil {
		panic("读取码表失败")
	}
	scan := bufio.NewScanner(rd)
	dict := make(map[rune][]string, 1e4)
out:
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			fmt.Println("这一行没有以\\t分割", scan.Text())
			continue
		}
		word, code := entry[0], entry[1]
		if len([]rune(word)) != 1 {
			// log.Println("不是单字，跳过", ws)
			continue
		}
		char := []rune(word)[0] // 是一个字
		// 去重，取较长的编码
		if _, ok := dict[char]; !ok {
			dict[char] = []string{code}
		} else {
			for i, v := range dict[char] {
				if strings.HasPrefix(v, code) { // 原来的编码更长
					continue out
				} else if strings.HasPrefix(code, v) { // 当前的编码更长
					dict[char][i] = code
					continue out
				}
			}
			dict[char] = append(dict[char], code)
		}
	}
	return dict
}

// 解析规则 2=AaAbBaBb,0=AaBaCaZa
func parseRule(s string) (map[int][]byte, []byte) {
	rule := make(map[int][]byte)

	// 整句规则
	ruleZ := make([]byte, 0, 1)
	if strings.HasSuffix(s, "...") {
		s = strings.ToLower(s)
		s = strings.TrimSuffix(s, "...")
		for i := range s {
			ruleZ = append(ruleZ, s[i]-'a')
		}
		return rule, ruleZ
	}

	s = strings.ReplaceAll(s, "，", ",")
	ruleSli := strings.Split(s, ",")

	for _, r := range ruleSli {
		// 对每条 2=AaAbBaBb 或 AABB
		tmp := strings.Split(r, "=")
		if len(tmp) != 2 {
			continue
		}
		idx, err := strconv.Atoi(tmp[0])
		if err != nil {
			panic("规则解析错误")
		}
		rule[idx] = partOfRule(tmp[1])
	}
	return rule, ruleZ
}

// AaAbBaBb 转为 0 0 0 1 1 0 1 1
// ABCC 转为 0 0 1 0 2 0 2 1
func partOfRule(s string) []byte {
	// 大写字母次数
	stat := make(map[byte]byte)
	ret := make([]byte, 0, 1)
	for i := 0; i < len(s); {
		if s[i] == 'Z' {
			ret = append(ret, 127)
		} else {
			ret = append(ret, s[i]-'A')
		}
		idx := stat[s[i]]
		stat[s[i]]++
		if i+1 < len(s) {
			if s[i+1] == 'z' {
				ret = append(ret, 127)
			} else if s[i+1] >= 'a' {
				ret = append(ret, s[i+1]-'a')
			} else {
				ret = append(ret, idx)
				i--
			}
			i += 2
		} else {
			ret = append(ret, idx)
			i++
		}
	}
	return ret
}
