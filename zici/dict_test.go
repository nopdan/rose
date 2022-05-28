package zici

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"testing"
)

func Test(t *testing.T) {

	// 哲哲豆词库 1w 多条
	format := "duoduo"
	dict := Parse(format, "test/duoduo.txt")
	write("out/"+format, dict)

	// 091 点儿词库
	format = "jidian_mb"
	dict = Parse(format, "test/jidian.mb")
	write("out/"+format, dict)

}

func TestOwn(t *testing.T) {

	// 哲哲豆词库
	format := "baidu_def"
	fp := "own/baidu.def"
	dict := Parse(format, fp)
	write(fp, dict)

}

func write(filename string, dict interface{}) {
	ioutil.WriteFile(fmt.Sprintf("%s.txt", filename), GenDuoduo(ToZcEntries(dict)), 0777)
}

// 用函数生成
func TestGen(t *testing.T) {
	// 091 点儿词库
	format := "jidian_mb"
	dict := Parse(format, "test/jidian.mb")

	ioutil.WriteFile("out/gen.txt", GenDuoduo(ToZcEntries(dict)), 0777)
}

// 用模版生成
func TestTmpl(t *testing.T) {
	// 091 点儿词库
	format := "jidian_mb"
	dict := Parse(format, "test/jidian.mb")

	f, _ := os.OpenFile("out/tmpl.txt", os.O_TRUNC|os.O_CREATE, 0777)
	defer f.Close()

	tmpl, _ := template.New("gen").Parse("{{ range . }}{{ .Word }}\t{{ .Code }}\n{{ end }}")
	err := tmpl.Execute(f, dict)
	if err != nil {
		panic(err)
	}
}
