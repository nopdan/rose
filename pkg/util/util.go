package util

import (
	"fmt"
	"io"
	"runtime"

	"golang.org/x/exp/constraints"
)

func Info[T constraints.Integer](r io.Reader, size T, info string) {
	tmp := ReadN(r, size)
	fmt.Printf("%s%s\n", info, DecodeMust(tmp, "UTF-16LE"))
}

func PrintHex(b []byte) {
	for _, v := range b {
		fmt.Printf("%02x ", v)
	}
	fmt.Println()
}

var LineBreak string

func init() {
	switch runtime.GOOS {
	case "windows":
		LineBreak = "\r\n"
	case "linux", "darwin":
		LineBreak = "\n"
	}
}
