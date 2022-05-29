package util

import "runtime"

var LineBreak = lineBreak()

func lineBreak() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	case "darwin":
		return "\r"
	default:
		return "\n"
	}
}
