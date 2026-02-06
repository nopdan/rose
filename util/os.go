package util

import (
	"runtime"
)

// LineBreak 根据操作系统返回合适的换行符
var LineBreak string

func init() {
	if runtime.GOOS == "windows" {
		LineBreak = "\r\n"
	} else {
		LineBreak = "\n"
	}
}
