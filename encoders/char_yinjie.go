package encoder

import (
	_ "embed"

	. "github.com/cxcn/dtool/utils"
)

//go:embed assets/char_yinjie.txt
var char_yinjie []byte

var CharYinjieMap = genCharYinjieMap()

// 生成单字音节表
func genCharYinjieMap() CharCodes {
	ret := ReadCharCodes(char_yinjie)
	// ascii
	var a byte = 32
	for ; a < 127; a++ {
		ret[rune(a)] = []string{string(a)}
	}
	return ret
}
