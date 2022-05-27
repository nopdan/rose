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
	write(format, dict)

	// // 哲哲豆词库
	// format = "baidu_def"
	// dict = Parse(format, "test/baidu.def")
	// write(format, dict)

	// 091 点儿词库
	format = "jidian_mb"
	dict = Parse(format, "test/jidian.mb")
	write(format, dict)

}

func write(filename string, dict interface{}) {
	// var buf bytes.Buffer
	// for _, v := range data {
	// 	buf.WriteString(v.Word + "\t" + v.Code + "\r\n")
	// }
	// ioutil.WriteFile(fmt.Sprintf("out/%s.txt", filename), buf.Bytes(), 0777)
	ioutil.WriteFile(fmt.Sprintf("out/%s.txt", filename), GenDuoduo(ToZcEntries(dict)), 0777)
}

func TestGen(t *testing.T) {
	// 091 点儿词库
	format := "jidian_mb"
	dict := Parse(format, "test/jidian.mb")

	ioutil.WriteFile("out/gen.txt", GenDuoduo(ToZcEntries(dict)), 0777)
}

func TestTmpl(t *testing.T) {
	// 091 点儿词库
	format := "jidian_mb"
	dict := Parse(format, "test/jidian.mb")

	f, _ := os.OpenFile("out/tmpl.txt", os.O_TRUNC|os.O_CREATE, 0777)
	tmpl, _ := template.New("gen").Parse("{{ range . }}{{ .Word }}\t{{ .Code }}\n{{ end }}")
	err := tmpl.Execute(f, dict)
	if err != nil {
		panic(err)
	}
	f.Close()
}
