package checker

import (
	"bufio"
	"bytes"
	"log"
	"strconv"
	"strings"

	"github.com/cxcn/dtool/pkg/table"
	"github.com/cxcn/dtool/pkg/util"
)

type Checker struct {
	Dict  map[rune][]string
	Rule  map[int][]byte
	RuleZ []byte
}

func NewChecker(path, rule string) *Checker {
	e := new(Checker)
	e.Dict = newDict(path)
	e.Rule, e.RuleZ = newRule(rule)
	return e
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
func newDict(path string) map[rune][]string {
	rd, err := util.Read(path)
	if err != nil {
		log.Fatal("读取码表失败")
	}
	scan := bufio.NewScanner(rd)
	d := make(map[rune][]string, 1e4)
out:
	for scan.Scan() {
		items := strings.Split(scan.Text(), "\t")
		if len(items) < 2 {
			log.Println("这一行没有以\\t分割", scan.Text())
			continue
		}
		ws, c := []rune(items[0]), items[1]
		if len(ws) != 1 {
			// log.Println("不是单字，跳过", ws)
			continue
		}
		w := ws[0] // 是一个字
		// 去重，取较长的编码
		if _, ok := d[w]; !ok {
			d[w] = []string{c}
		} else {
			for i, v := range d[w] {
				if strings.HasPrefix(v, c) { // 原来的编码更长
					continue out
				} else if strings.HasPrefix(c, v) { // 当前的编码更长
					d[w][i] = c
					continue out
				}
			}
			d[w] = append(d[w], c)
		}
	}
	return d
}

// 处理规则 2=AaAbBaBb,0=AaBaCaZa
func newRule(rule string) (map[int][]byte, []byte) {
	retMap := make(map[int][]byte)
	retSli := make([]byte, 0, 1)

	if strings.HasSuffix(rule, "...") {
		rule = strings.ToLower(rule)
		rule = strings.TrimSuffix(rule, "...")
		for i := range rule {
			retSli = append(retSli, rule[i]-'a')
		}
		return retMap, retSli
	}

	rule = strings.ReplaceAll(rule, "，", ",")
	sliRule := strings.Split(rule, ",")

	for _, r := range sliRule {
		// 对每条 2=AaAbBaBb 或 AABB
		tmp := strings.Split(r, "=")
		if len(tmp) != 2 {
			continue
		}
		idx, err := strconv.Atoi(tmp[0])
		if err != nil {
			panic("规则解析错误")
		}
		retMap[idx] = fRule(tmp[1])
	}
	return retMap, retSli
}

// AaAbBaBb 转为 0 0 0 1 1 0 1 1
// ABCC 转为 0 0 1 0 2 0 2 1
func fRule(s string) []byte {
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
