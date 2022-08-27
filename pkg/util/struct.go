package util

// 单字码表
type CharCodes map[rune][]string

// 多多形式码表
type WcTable []WordCode

// 词，编码
type WordCode struct {
	Word string
	Code string
}

// 极点形式码表
type CwsTable []CodeWords

// 一码多词
type CodeWords struct {
	Code  string
	Words []string
}

// 拼音词库
type WpfDict []WordPyFreq

// 词，拼音，词频
type WordPyFreq struct {
	Word   string
	Pinyin []string
	Freq   int
}

// // 不带词频
// type WpDict []WordPinyin

// // 词，拼音
// type WordPinyin struct {
// 	Word   string
// 	Pinyin []string
// }
