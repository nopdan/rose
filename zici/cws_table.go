package zici

import (
	"sort"

	. "github.com/cxcn/dtool/utils"
)

// 转为 极点形式，一码多词
func ToCwsTable(wct WcTable) CwsTable {
	codeMap := make(map[string][]string)
	for _, wc := range wct {
		if _, ok := codeMap[wc.Code]; !ok {
			codeMap[wc.Code] = []string{wc.Word}
			continue
		}
		codeMap[wc.Code] = append(codeMap[wc.Code], wc.Word)
	}
	ret := make(CwsTable, 0, len(codeMap))
	for k, v := range codeMap {
		ret = append(ret, CodeWords{k, v})
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Code < ret[j].Code
	})
	return ret
}
