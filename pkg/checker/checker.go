package checker

import (
	"bufio"
	"log"
	"strconv"
	"strings"

	"github.com/cxcn/dtool/pkg/table"
	"github.com/cxcn/dtool/pkg/util"
)

type Checker struct {
	Dict map[rune][]string
	Rule map[int][]int
}

func NewChecker(path, rule string) *Checker {
	e := new(Checker)
	e.Dict = newDict(path)
	e.Rule = newRule(rule)
	return e
}

func (c *Checker) Check(table table.Table) string {
	var sb strings.Builder
	sb.WriteString("词\t编码\t可能的编码\n")
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
			sb.WriteString(word + "\t" + code + "\t")
			sb.WriteString(strings.Join(codes, " "))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
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
func newRule(rule string) map[int][]int {
	rule = strings.ReplaceAll(rule, "，", ",")
	sliRule := strings.Split(rule, ",")
	// 把 AaAbBaBb 转为了 1 1 1 2 2 1 2 2
	f := func(s string) []int {
		s = strings.TrimSpace(s)
		// 65 97
		ret := make([]int, 0, 10)
		for i := range s {
			id := 0
			if s[i] == 90 || s[i] == 122 {
				// Z and z 用 0 表示
			} else if s[i] >= 97 { // A-Z
				id = int(s[i]) - 96
			} else { // a-z
				id = int(s[i]) - 64
			}
			ret = append(ret, id)
		}
		return ret
	}

	ret := make(map[int][]int)
	for _, r := range sliRule {
		tmp := strings.Split(r, "=")
		if len(tmp) != 2 {
			continue
		}
		id, err := strconv.Atoi(tmp[0])
		if err != nil {
			log.Fatal("规则解析错误")
		}
		ret[id] = f(tmp[1])
	}
	return ret
}
