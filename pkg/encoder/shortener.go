package encoder

import (
	"strconv"
	"strings"

	. "github.com/cxcn/dtool/pkg/util"
)

// rule 1:0,2:3,3:2,6:n 默认1，n 无限
func Shorten(wct *WcTable, rule string) {
	rl := handleRule(rule)
	countMap := make(map[string]int)
	for i := 0; i < len(*wct); i++ {
		wc := (*wct)[i]
		for j := 1; j <= len(wc.Code); j++ {
			curr := string(wc.Code[:j])
			count := countMap[curr]
			if count < rl[j] {
				wc.Code = curr
				(*wct)[i] = wc
				countMap[curr]++
				break
			}
		}
	}
	// fmt.Println(countMap)
}

// [0,3,2,1,1,1e5]
func handleRule(rule string) []int {
	ret := make([]int, 0)
	rule = strings.ReplaceAll(rule, " ", "")
	rule = strings.ReplaceAll(rule, "，", ",")
	r := strings.Split(rule, ",")
	for _, v := range r {
		tmp := strings.Split(v, ":")
		if len(tmp) != 2 {
			continue
		}
		pos, _ := strconv.Atoi(tmp[0])
		if pos < 1 {
			continue
		}
		var val int
		if tmp[1] == "n" {
			val = 1e5
		} else {
			val, _ = strconv.Atoi(tmp[1])
		}
		setVal(&ret, pos, val)
	}
	// fmt.Println(ret)
	return ret
}

func setVal(sli *[]int, pos int, val int) {
	for pos > len(*sli)-1 {
		*sli = append(*sli, 1)
	}
	(*sli)[pos] = val
}
