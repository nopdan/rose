package encoder

import (
	"fmt"
	"testing"
)

func TestGetPinyin(t *testing.T) {
	fmt.Println(GetPinyin("会计师"))
	// ！！ 顺序最长匹配，人参 / 加
	fmt.Println(GetPinyin("一个人参加了会议"))
}
