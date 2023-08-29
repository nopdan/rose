package core

import (
	"fmt"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	c := &Config{
		IName:   "sample/words.txt",
		IFormat: "words",
		OFormat: "rime",
		OName:   "test/to_rime.txt",
	}
	c.Marshal()
}

func TestFormatList(t *testing.T) {
	// 检查 id 是否有重复
	formatSet := make(map[string]struct{})
	for _, f := range FormatList {
		ids := strings.Split(f.ID, ",")
		for _, id := range ids {
			if _, ok := formatSet[id]; ok {
				t.Fatalf("id %s 重复", id)
			}
			formatSet[id] = struct{}{}
		}
	}
	fmt.Println("id 无重复")
	PrintFormatList()
}
