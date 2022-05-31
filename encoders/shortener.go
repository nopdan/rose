package encoder

import (
	"strconv"
	"strings"
)

type Table []*struct {
	Word string
	Code string
}

// rule 1:0,2:3,3:2,6:n 默认1，n 无限
func (t Table) Shorten(rule string) {
	rl := handleRule(rule)
	countMap := make(map[string]int)
	for _, v := range t {
		for i := 1; i <= len(v.Code); i++ {
			curr := string(v.Code[:i])
			count := countMap[curr]
			if count < rl[i] {
				v.Code = curr
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
	r := strings.Split(rule, ",")
	for _, v := range r {
		v = strings.TrimSpace(v)
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
