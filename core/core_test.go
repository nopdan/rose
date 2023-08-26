package core

import (
	"testing"
)

func TestMain(t *testing.T) {
	c := &Config{
		Name:      "sample/words.txt",
		InFormat:  "words",
		OutFormat: "rime",
		OutName:   "test/to_rime.txt",
	}
	c.Marshal()
}
