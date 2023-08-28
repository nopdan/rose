package encoder

import (
	"fmt"
	"testing"
)

func TestMatch(t *testing.T) {

	enc := NewPinyin()

	fmt.Println(enc.Encode("会计师"))

	fmt.Println(enc.Encode("一个人参加了会议"))

	fmt.Println(enc.Encode("α"))
}
