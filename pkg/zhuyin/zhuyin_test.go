package zhuyin

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	fmt.Println(Get("会计师"))
	// ！！ 顺序最长匹配，人参 / 加
	fmt.Println(Get("一个人参加了会议"))

	fmt.Println(Get("α"))
}
