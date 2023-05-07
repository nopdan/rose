package pinyin

import (
	"fmt"
	"testing"
)

func TestMatch(t *testing.T) {
	fmt.Println(Match("会计师"))

	fmt.Println(Match("一个人参加了会议"))

	fmt.Println(Match("α"))
}
