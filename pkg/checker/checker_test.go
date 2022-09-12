package checker

import (
	"fmt"
	"testing"

	"github.com/cxcn/dtool/pkg/table"
)

func TestRule(t *testing.T) {
	fmt.Println(newRule("2=AaAbBaBb,3=AaAbBaCa,0=AaBaCaZa"))
	fmt.Println(newRule("2=AaAbBaBb,3=AaBaCaCb,0=AaBaCaZa"))
	fmt.Println(newRule("2=AaAbBaBbAcBc,3=AaBaCaAcBcCc,0=AaBaCaZaAcBc"))
	fmt.Println(newRule("2=AABB,3=ABCC,0=ABCZ"))
	fmt.Println(newRule("ab..."))
}

func TestChecker(t *testing.T) {
	rule := "2=AABB,3=AABC,0=ABCZ"
	path := "./own/test.txt"
	c := NewChecker(path, rule)
	table := table.Parse("duoduo", path)
	fmt.Println(string(c.Check(table)))

	tmp := c.Encode("温柔\n没有人\n好不容易\n对外贸易法")
	for word, codes := range tmp {
		fmt.Println(word, "\n", codes)
	}
}
