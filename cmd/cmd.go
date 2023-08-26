package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nopdan/rose/core"
)

func Cmd() {
	switch len(os.Args) {
	case 1:
		ask()
		return
	case 2:
		switch os.Args[1] {
		case "-v", "version":
			fmt.Println("蔷薇词库转换v1.1.3\nhttps://github.com/nopdan/rose")
			return
		case "-h", "help":
			fmt.Printf("Usage: .\\rose.exe [input] [input_format]:[output_format] [output]\n")
			fmt.Printf("Example: .\\rose.exe .\\sogou.scel scel:rime rime.dict.yaml\n")
			return
		}
	}

	if len(os.Args) >= 3 {
		goto NORMAL
	}
	wrong()
	return

NORMAL:
	c := &core.Config{}

	c.Name = os.Args[1]
	fm := strings.Split(os.Args[2], ":")
	if len(fm) == 2 {
		c.InFormat = fm[0]
		c.OutFormat = fm[1]
	}
	if len(os.Args) > 3 {
		c.OutName = os.Args[3]
	}
	// fmt.Println(input, output, input_format, output_format)
	c.Marshal()
}

func wrong() {
	fmt.Println("输入参数有误")
}

// 交互式
func ask() {
	askOne := func(hint string) string {
		fmt.Printf("%s\n> ", hint)
		reader := bufio.NewReader(os.Stdin)
		var value string
		value, _ = reader.ReadString('\n')
		value = strings.ReplaceAll(value, "\r", "")
		value = strings.ReplaceAll(value, "\n", "")
		last := value[len(value)-1]
		if len(value) > 3 {
			if last == '"' && value[0] == '"' {
				value = value[1 : len(value)-1]
			} else if len(value) > 4 && last == '\'' && value[:3] == "& '" {
				value = value[3 : len(value)-1]
			}
		}

		fmt.Println()
		return value
	}
	c := &core.Config{}

	c.Name = askOne("输入词库：")
	c.InFormat = askOne("词库格式：")
	c.OutName = askOne("输出词库：")
	c.OutFormat = askOne("词库格式：")
	c.Marshal()

	// fmt.Scanln()
}
