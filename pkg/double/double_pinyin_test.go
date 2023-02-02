package double

import (
	"fmt"
	"testing"

	"github.com/imetool/dtool/pkg/pinyin"
)

func TestMapping(t *testing.T) {
	m := newMapping("../../assets/双拼映射表/星辰双拼.ini", AABC)
	fmt.Println(m)
}

func TestToDoublePinyin(t *testing.T) {
	dict := pinyin.Parse("sogou_bin", "../pinyin/test/sogou_bak.bin")
	table := ToDoublePinyin(dict, "test/双拼映射表.ini", AABC)
	fmt.Println(table)
}
