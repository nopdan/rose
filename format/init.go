package format

import (
	"github.com/nopdan/rose/format/baidu_bak"
	"github.com/nopdan/rose/format/baidu_bdict"
	"github.com/nopdan/rose/format/baidu_def"
	"github.com/nopdan/rose/format/custom_text"
	"github.com/nopdan/rose/format/fcitx4_mb"
	"github.com/nopdan/rose/format/jidian"
	"github.com/nopdan/rose/format/jidian_mb"
	"github.com/nopdan/rose/format/kafan_pinyin_bak"
	"github.com/nopdan/rose/format/kafan_wubi_bak"
	"github.com/nopdan/rose/format/mspy_udl"
	"github.com/nopdan/rose/format/msudp"
	"github.com/nopdan/rose/format/mswb_lex"
	"github.com/nopdan/rose/format/pinyinjiajia"
	"github.com/nopdan/rose/format/qq_qpyd"
	"github.com/nopdan/rose/format/rime"
	"github.com/nopdan/rose/format/sogou_bak"
	"github.com/nopdan/rose/format/sogou_scel"
	"github.com/nopdan/rose/format/ziguang_uwl"
)

// 初始化时注册所有格式
func init() {
	// 二进制拼音
	RegisterFormat(sogou_scel.New())
	RegisterFormat(sogou_scel.NewQcel())
	RegisterFormat(sogou_bak.NewSogouBak())
	RegisterFormat(baidu_bdict.New())
	RegisterFormat(baidu_bdict.NewBcd())
	RegisterFormat(baidu_bak.New())
	RegisterFormat(qq_qpyd.New())
	RegisterFormat(mspy_udl.New())
	RegisterFormat(ziguang_uwl.New())
	RegisterFormat(kafan_pinyin_bak.New())

	// 纯文本拼音
	RegisterFormat(pinyinjiajia.New())
	RegisterFormat(rime.New())
	RegisterFormat(custom_text.NewSogouPinyin())
	RegisterFormat(custom_text.NewBaiduPinyin())
	RegisterFormat(custom_text.NewQQPinyin())

	// 二进制五笔
	RegisterFormat(baidu_def.New())
	RegisterFormat(msudp.New())
	RegisterFormat(mswb_lex.New())
	RegisterFormat(jidian_mb.New())
	RegisterFormat(fcitx4_mb.New())
	RegisterFormat(kafan_wubi_bak.New())

	// 纯文本五笔
	RegisterFormat(jidian.New())
	RegisterFormat(custom_text.NewDuoduoWubi())
	RegisterFormat(custom_text.NewBaiduWubi())
	RegisterFormat(custom_text.NewUserPhrase())

	// 纯词组
	RegisterFormat(custom_text.NewWords())
}
