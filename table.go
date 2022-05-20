package dtool

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

// 词库 code: []word
type Table map[string][]string
type sTable struct {
	code  string
	words []string
}

// 初始化词库
func NewTable() Table {
	return make(Table)
}

func (t Table) Check(d Dict, rules map[int][]int) {
	// 遍历整个码表
	for code, words := range t {
		// 对每个词
		for _, word := range words {
			// 排除单字
			if len([]rune(word)) == 1 {
				continue
			}
			flag := false
			// 根据规则推测可能的编码
			codes := d.encodeWord(word, rules)
			// 对每个编码
			for _, c := range codes {
				if strings.HasPrefix(c, code) {
					flag = true
					break
				}
			}
			if !flag && len(codes) != 0 {
				fmt.Println("!check failed: ", word, code, codes)
			}
		}
	}
}

// 保存
func (t Table) Save(filename string) {
	res := make([]sTable, 0, len(t))
	for code, words := range t {
		res = append(res, sTable{code, words})
	}
	// 排序
	sort.Slice(res, func(i, j int) bool {
		return res[i].code < res[j].code
	})

	var b bytes.Buffer
	for i := 0; i < len(res); i++ {
		for j := 0; j < len(res[i].words); j++ {
			b.WriteString(res[i].words[j])
			b.WriteByte('\t')
			b.WriteString(res[i].code)
			b.WriteByte('\n')
		}
	}
	ioutil.WriteFile(filename, b.Bytes(), 0777)
}

// 合并码表
func (t Table) Merge(m map[string][]string) {
	for word, codes := range m {
		for _, code := range codes {
			t[code] = append(t[code], word)
			t[code] = rmDup(t[code])
		}
	}
}

// 删除切片重复元素
func rmDup(arr []string) []string {
	m := make(map[string]struct{}, len(arr))
	ret := make([]string, 0, len(arr))
	for _, v := range arr {
		// 没有重复就添加进切片
		if _, ok := m[v]; !ok {
			ret = append(ret, v)
			m[v] = struct{}{}
		}
	}
	return ret
}

// 读取码表，flag: 字在后
func (t Table) Read(s string, flag bool) {
	f, err := os.Open(s)
	rd, err := ReadFile(f)
	if err != nil {
		log.Fatal("读取码表失败")
	}
	scan := bufio.NewScanner(rd)

	for scan.Scan() {
		items := strings.Split(scan.Text(), "\t")
		if len(items) < 2 {
			log.Println("这一行没有以\\t分割", scan.Text())
			continue
		}
		if flag { // 编码在前
			t[items[0]] = append(t[items[0]], items[1])
		} else { // 编码在后
			t[items[1]] = append(t[items[1]], items[0])
		}
	}
}
