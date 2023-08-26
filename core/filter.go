package core

type FilterConfig struct {
	WordLenMin int // 过滤小于这个长度的词
	WordLenMax int // 过滤大于这个长度的词
	FreqLenMin int // 过滤词频小于这个值的词
	FreqLenMax int // 过滤词频大于这个值的词
	// 根据词过滤
	HasAlphabet bool // 过滤含英文的词
	HasNumber   bool // 过滤含数字的词
}