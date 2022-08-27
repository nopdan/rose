package main

import (
	"fmt"
	"os"

	"github.com/cxcn/dtool/pkg/dtool"
	"github.com/jessevdk/go-flags"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("\nPress the enter key to exit...")
		fmt.Scanln()
	} else {
		cli()
	}
}

func cli() {
	var opts struct {
		Input   string `short:"i" description:"string\t输入词库"`
		IFormat string `short:"f" description:"string\t输入格式"`
		Output  string `short:"o" description:"string\t输出词库"`
		OFormat string `short:"w" description:"string\t输出格式"`
	}

	flags.Parse(&opts)
	if opts.Input == "" {
		return
	}
	if opts.IFormat == "" {
		opts.IFormat = "baidu"
	}
	if opts.OFormat == "" {
		opts.OFormat = "baidu"
	}
	if opts.Output == "" {
		opts.Output = opts.Input + opts.OFormat + ".txt"
	}

	pes := dtool.PinyinParse(opts.IFormat, opts.Input)
	data := dtool.PinyinGen(opts.OFormat, pes)
	os.WriteFile(opts.Output, data, 0777)
}
