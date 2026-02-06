package encoder

import (
	"bufio"
	"bytes"
	_ "embed"
	"strings"

	"github.com/nopdan/rose/model"
)

//go:embed CJK.txt
var cjkData []byte

// WubiEncoder 五笔编码器
type WubiEncoder struct {
	*BaseEncoder
	charMap map[rune]string // 字符到编码的映射
	schema  string
	isAABC  bool
}

// NewWubiEncoder 创建五笔编码器
func NewWubiEncoder(schema string, data []byte, isAABC bool) *WubiEncoder {
	encoder := &WubiEncoder{
		BaseEncoder: &BaseEncoder{},
		charMap:     make(map[rune]string),
		schema:      schema,
		isAABC:      isAABC,
	}

	// 加载编码数据
	if schema == "custom" && data != nil {
		encoder.initCustomData(data)
	} else {
		encoder.initWubiData()
	}
	return encoder
}

// initCustomData 加载自定义五笔编码数据
func (e *WubiEncoder) initCustomData(data []byte) {
	r := bytes.NewReader(data)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) < 2 {
			continue
		}
		char := []rune(fields[0])
		if len(char) != 1 {
			continue
		}
		code := fields[1]
		e.charMap[char[0]] = code
	}
}

// initWubiData 加载五笔编码数据
func (e *WubiEncoder) initWubiData() {
	r := bytes.NewReader(cjkData)
	scanner := bufio.NewScanner(r)

	// 根据schema确定使用第几列
	var colIndex int
	switch e.schema {
	case "wubi86", "86":
		colIndex = 2
	case "wubi98", "98":
		colIndex = 3
	case "wubi06", "06":
		colIndex = 4
	default:
		colIndex = 2 // 默认使用86版
	}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) <= colIndex || len(fields) < 2 {
			continue
		}

		// 获取字符和编码
		chars := []rune(fields[1])
		if len(chars) == 1 {
			char := chars[0]
			code := fields[colIndex]
			if code != "" {
				e.charMap[char] = code
			}
		}
	}
}

// Encode 将Entry转换为五笔编码（原地修改）
func (e *WubiEncoder) Encode(entry *model.Entry) {
	// 获取五笔编码
	wubiCode := e.encodeWord(entry.Word)

	// 如果编码为空，保持原来的编码
	if wubiCode == "" {
		return
	}

	entry.Code = model.NewSimpleCode(wubiCode)
	entry.CodeType = model.CodeTypeWubi
}

// encodeWord 将汉字词组转换为五笔编码
func (e *WubiEncoder) encodeWord(word string) string {
	wordRunes := []rune(word)
	wordLen := len(wordRunes)
	if wordLen == 0 {
		return ""
	}

	var code string
	switch wordLen {
	case 1:
		// 单字
		if c, ok := e.charMap[wordRunes[0]]; ok {
			code = c
		}
	case 2:
		// 两字词
		a := e.getCharCode(wordRunes[0])
		b := e.getCharCode(wordRunes[1])
		code = e.cutCode(a, 2) + e.cutCode(b, 2)
	case 3:
		// 三字词
		a := e.getCharCode(wordRunes[0])
		b := e.getCharCode(wordRunes[1])
		c := e.getCharCode(wordRunes[2])

		if e.isAABC {
			code = e.cutCode(a, 2) + e.cutCode(b, 1) + e.cutCode(c, 1)
		} else {
			code = e.cutCode(a, 1) + e.cutCode(b, 1) + e.cutCode(c, 2)
		}
	default:
		// 四字及以上词
		a := e.getCharCode(wordRunes[0])
		b := e.getCharCode(wordRunes[1])
		c := e.getCharCode(wordRunes[2])
		z := e.getCharCode(wordRunes[wordLen-1])
		code = e.cutCode(a, 1) + e.cutCode(b, 1) + e.cutCode(c, 1) + e.cutCode(z, 1)
	}

	return code
}

// getCharCode 获取字符的编码
func (e *WubiEncoder) getCharCode(char rune) string {
	if code, ok := e.charMap[char]; ok {
		return code
	}
	return ""
}

// cutCode 截取编码的指定长度
func (e *WubiEncoder) cutCode(code string, length int) string {
	if len(code) < length {
		return code
	}
	return code[:length]
}

// EncodeBatch 批量编码（原地修改）
func (e *WubiEncoder) EncodeBatch(entries []*model.Entry) {
	for _, entry := range entries {
		e.Encode(entry)
	}
}
