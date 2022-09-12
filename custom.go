package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/cxcn/dtool/pkg/checker"
	"github.com/cxcn/dtool/pkg/encoder"
	"github.com/cxcn/dtool/pkg/pinyin"
	"github.com/cxcn/dtool/pkg/table"
	"github.com/cxcn/dtool/pkg/util"
)

// 选择文件
func (a *App) SelectFile() string {
	opts := runtime.OpenDialogOptions{}
	ret, _ := runtime.OpenFileDialog(a.ctx, opts)
	return ret
}

func (a *App) getSavePath(base string) (string, runtime.MessageDialogOptions) {
	// 选择保存位置
	opts := runtime.SaveDialogOptions{
		DefaultDirectory: filepath.Dir(base),
	}
	ret, _ := runtime.SaveFileDialog(a.ctx, opts)

	mdo := runtime.MessageDialogOptions{
		Type:    "Ok",
		Title:   "Success!",
		Message: "保存成功！",
		Buttons: []string{"确认"},
	}
	return ret, mdo
}

func (a *App) Convert(input, iformat, oformat string, isPinyin bool) {
	ret, mdo := a.getSavePath(input)
	// 没有选
	if ret == "" {
		return
	}
	var data []byte
	if isPinyin {
		dict := pinyin.Parse(iformat, input)
		data = pinyin.Generate(oformat, dict)
	} else {
		dict := table.Parse(iformat, input)
		data = table.Generate(oformat, dict)
	}
	os.WriteFile(ret, data, 0666)
	runtime.MessageDialog(a.ctx, mdo)
}

func (a *App) Shorten(input, rule string) {
	ret, mdo := a.getSavePath(input)
	// 没有选
	if ret == "" {
		return
	}
	wct := table.Parse("duoduo", input)
	encoder.Shorten(&wct, rule)
	data := table.Generate("duoduo", wct)
	os.WriteFile(ret, data, 0666)
	runtime.MessageDialog(a.ctx, mdo)
}

func (a *App) Encode(charPath, dictPath, encRule string, isCheck bool) {
	savePath, mdo := a.getSavePath(charPath)
	// 没有选
	if savePath == "" {
		return
	}
	ck := checker.NewChecker(charPath, encRule)
	if isCheck {
		tb := table.Parse("duoduo", dictPath)
		data := ck.Check(tb)
		os.WriteFile(savePath, []byte(data), 0666)
	} else {
		rd, _ := util.Read(dictPath)
		b, _ := io.ReadAll(rd)
		data := ck.Encode(string(b))
		var buf bytes.Buffer
		for i := range data {
			for _, code := range data[i].Codes {
				buf.WriteString(data[i].Word)
				buf.WriteByte('\t')
				buf.WriteString(code)
				buf.WriteString(util.LineBreak)
			}
		}
		os.WriteFile(savePath, buf.Bytes(), 0666)
	}
	runtime.MessageDialog(a.ctx, mdo)
}
