package checker

import (
	"fmt"
	"testing"

	"github.com/cxcn/dtool/pkg/table"
)

func TestChecker(t *testing.T) {
	// rule := "2=AaAbBaBb\r\n3=AaBaCaCb\r\n0=AaBaCaZa"
	rule := "2=AaAbBaBb\r\n3=AaAbBaCa\r\n0=AaBaCaZa"
	// rule := "2=AaAbBaBbAcBc\r\n3=AaBaCaAcBcCc\r\n0=AaBaCaZaAcBc"
	path := "./own/test.txt"
	c := NewChecker(path, rule)
	table := table.Parse("duoduo", path)
	c.Check(table)
	// for k, v := range c.Dict {
	// 	fmt.Println(string(k), v)
	// }

	c.EncodeWord("你们")
	c.EncodeWord("在意")
	c.EncodeWord("参加")
	c.EncodeWord("行人")
	c.EncodeWord("adg")
	tmp := c.Encode("温柔\n没有人\n好不容易\n对外贸易法")
	for word, codes := range tmp {
		fmt.Println(word, "\n", codes)
	}
}
