package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	server "github.com/nopdan/rose/frontend"
	"github.com/nopdan/rose/pkg/core"
)

func Cmd() {
	// 双击打开默认启动服务
	switch len(os.Args) {
	case 1:
		server.Serve(7800)
		return
	}
	switch os.Args[1] {
	case "list":
		core.PrintFormatList()
	case "help", "-h":
		help()
	case "version", "-v":
		fmt.Printf("Go Version: %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Rose Version: v1.2.1\n")
		fmt.Printf("Author: nopdan <me@nopdan.com>\n")
	case "server":
		port := 7800
		if len(os.Args) > 2 {
			p := os.Args[2]
			if strings.HasPrefix(p, "-p:") {
				p = p[3:]
				tmp, err := strconv.Atoi(p)
				if err != nil {
					port = tmp
				}
			}
		}
		server.Serve(port)
	default:
		root()
	}
}

func help() {
	fmt.Println("Root Command:")
	fmt.Printf("    Usage: rose [输入文件] [输入格式]:[输出格式] [保存文件名]\n")
	fmt.Printf("    Example: rose sogou.scel scel:rime rime.dict.yaml\n")
	fmt.Println()
	fmt.Println("Sub Commands:")
	fmt.Println("      list      列出所有支持的格式")
	fmt.Println("      server    启动服务  -p:[port] 指定端口(默认7800)")
	fmt.Println("  -h, help      帮助")
	fmt.Println("  -v, version   版本")
}

func root() {
	if len(os.Args) <= 2 {
		fmt.Println("缺少必要参数")
		return
	}
	c := &core.Config{}
	c.IName = os.Args[1]
	fm := strings.Split(os.Args[2], ":")
	if len(fm) == 2 {
		c.IFormat = fm[0]
		c.OFormat = fm[1]
	} else {
		fmt.Println("输入输出格式参数解析错误")
		return
	}
	if len(os.Args) == 4 {
		c.OName = os.Args[3]
	}
	d := c.Marshal()
	c.Save(d)
}
