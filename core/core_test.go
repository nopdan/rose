package core

import (
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
	PrintFormatList()
}
