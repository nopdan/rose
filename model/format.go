package model

import (
	"fmt"
	"io"
	"os"
)

// FormatType 表示格式类型
type FormatType int

const (
	FormatTypePinyin FormatType = iota + 1 // 拼音类格式
	FormatTypeWubi                         // 五笔类格式
	FormatTypeWords                        // 纯词组格式
	FormatTypeCustom                       // 自定义格式
)

func (t FormatType) String() string {
	switch t {
	case FormatTypePinyin:
		return "拼音"
	case FormatTypeWubi:
		return "五笔"
	case FormatTypeWords:
		return "词组"
	case FormatTypeCustom:
		return "自定义"
	default:
		return "未知"
	}
}

// Format 格式元信息接口（仅描述）
type Format interface {
	// 获取格式信息
	Info() *BaseFormat
}

// Importer 导入接口
type Importer interface {
	// 导入：统一输入来源（本地文件/内存数据）
	Import(src Source) ([]*Entry, error)
}

// Exporter 导出接口
type Exporter interface {
	// 导出：将Entry列表输出到Writer
	Export(entries []*Entry, w io.Writer) error
}

// LogLevel 日志级别
// 0: 不打印任何信息
// 1: 打印解析或生成的基础信息
// 2: 测试用途，可输出解析或生成后的文件内容
type LogLevel int

const (
	LogNone LogLevel = iota
	LogInfo
	LogDebug
)

// BaseFormat 格式信息
type BaseFormat struct {
	// 格式唯一标识符，如 "sogou_scel"
	ID string
	// 显示名称，如 "搜狗细胞词库"
	Name string
	// 格式类型
	Type FormatType
	// 是否为二进制格式
	IsBinary bool
	// 支持的文件扩展名
	Extension string
	// 格式描述
	Description string

	// 日志级别
	LogLevel LogLevel
	// 日志文件
	LogFile *os.File
}

func (f *BaseFormat) Logf(level LogLevel, s string, v ...any) {
	switch f.LogLevel {
	case LogInfo:
		if level == LogInfo {
			fmt.Printf(s, v...)
		}
	case LogDebug:
		if f.LogFile == nil {
			f.LogFile, _ = os.Create("test_import.txt")
		}
		if f.LogFile != nil {
			fmt.Fprintf(f.LogFile, s, v...)
		}
	}
}

func (f *BaseFormat) Infof(s string, v ...any) {
	f.Logf(LogInfo, s, v...)
}

func (f *BaseFormat) Debugf(s string, v ...any) {
	f.Logf(LogDebug, s, v...)
}

func (f *BaseFormat) Info() *BaseFormat {
	return f
}
