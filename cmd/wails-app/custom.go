package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/cxcn/dtool/pinyin"
	"github.com/cxcn/dtool/zici"
)

// 选择文件
func (a *App) SelectFile() string {
	opts := runtime.OpenDialogOptions{}
	ret, _ := runtime.OpenFileDialog(a.ctx, opts)
	return ret
}

func (a *App) ConvDict(input, iformat, oformat string, isPy bool) {
	// 选择保存位置
	opts := runtime.SaveDialogOptions{
		DefaultDirectory: filepath.Dir(input),
	}
	ret, _ := runtime.SaveFileDialog(a.ctx, opts)

	mdo := runtime.MessageDialogOptions{
		Type:    "Ok",
		Title:   "DTool",
		Message: "保存成功！",
		Buttons: []string{"确认"},
	}
	// 没有选
	if ret == "" {
		return
	}
	var data []byte
	if isPy {
		data = ConvPyDict(input, iformat, oformat)
	} else {
		data = ConvZcDict(input, iformat, oformat)
	}
	ioutil.WriteFile(ret, data, 0777)
	runtime.MessageDialog(a.ctx, mdo)
}

// 转换拼音词库
func ConvPyDict(input, iformat, oformat string) []byte {
	pes := pinyin.Parse(iformat, input)
	data := pinyin.Gen(oformat, pes)
	return data
}

// 转换字词码表
func ConvZcDict(input, iformat, oformat string) []byte {
	pes := zici.Parse(iformat, input)
	data := zici.Gen(oformat, pes)
	return data
}
