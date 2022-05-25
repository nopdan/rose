package zici

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test(t *testing.T) {
	f, err := os.Open("jidian.mb")
	if err != nil {
		log.Panic(err)
	}
	d := ParseJidianMb(f)
	// fmt.Println(d)
	b := Gen("duoduo", d)
	ioutil.WriteFile("res.txt", b, 0777)
}
