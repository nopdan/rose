package encoder

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

// 原始词库来源：
// 现代汉语常用词表（2008年李行健课题组）.txt
// 深蓝 WordPinyin.txt
// 去除标点
// 从原始词库生成 word_pinyin.txt
func TestWordPinyin(t *testing.T) {
	f, _ := os.Open("own/src_word_pinyin.txt")
	defer f.Close()
	scan := bufio.NewScanner(f)

	var buf bytes.Buffer
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		word, codes := entry[0], entry[1]
		if codes != strings.Join(GetPyByChar(word), "'") {
			buf.WriteString(word)
			buf.WriteByte('\t')
			buf.WriteString(codes)
			buf.WriteString("\r\n")
		}
	}
	os.WriteFile("assets/word_pinyin.txt", buf.Bytes(), 0666)
}
