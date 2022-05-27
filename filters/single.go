package filter

type entry struct {
	word  string   // 词
	code  string   // 编码，定长方案的
	codes []string // 编码，拼音切片
	freq  int      // 词频
}

// 传入词条，回调函数，是否保留某个词条
func Filter(d []entry, f func(entry) bool) []entry {
	ret := make([]entry, 0, len(d))
	for _, e := range d {
		if f(e) {
			ret = append(ret, e)
		}
	}
	return ret
}

// 单字滤镜（只保留单字）
func Single(e entry) bool {
	// 单字，保留
	return len([]rune(e.word)) == 1
}
