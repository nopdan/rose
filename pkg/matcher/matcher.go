package matcher

type Matcher interface {
	// 插入一个词条 word code pos
	Insert(string, string, int)
	// 构建
	Build()
	// 匹配下一个词，返回匹配到的词长，编码和候选位置
	Match([]rune) (int, string, int)
}
