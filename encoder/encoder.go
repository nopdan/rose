package encoder

import "github.com/nopdan/rose/model"

// Encoder 编码器接口
type Encoder interface {
	Encode(entry *model.Entry)
	EncodeBatch(entries []*model.Entry)
}

// NullEncoder 空编码器（不进行转换）
type NullEncoder struct{}

func (e *NullEncoder) Encode(entry *model.Entry) {
	// 不进行转换
}

func (e *NullEncoder) EncodeBatch(entries []*model.Entry) {
	// 不进行转换
}

// BaseEncoder 基础编码器实现
type BaseEncoder struct{}

// Encode 基础实现，子类应该覆盖此方法
func (e *BaseEncoder) Encode(entry *model.Entry) {
	// 默认不转换
}

// EncodeBatch 批量编码
func (e *BaseEncoder) EncodeBatch(entries []*model.Entry) {
	for _, entry := range entries {
		e.Encode(entry)
	}
}

// NewEncoder 创建编码器的工厂函数
// params 支持:
//   - "schema": string  五笔方案 "86"/"98"/"06"/"custom"
//   - "useAABC": bool   组词规则
//   - "codeTableData": []byte  自定义码表数据
func NewEncoder(encoderType string, params map[string]any) Encoder {
	switch encoderType {
	case "none":
		return &NullEncoder{}
	case "pinyin", "py":
		return NewPinyinEncoder()
	case "wubi":
		schema := "wubi86"
		useAABC := false
		var customData []byte
		if params != nil {
			if s, ok := params["schema"].(string); ok && s != "" {
				schema = s
			}
			if aabc, ok := params["useAABC"].(bool); ok {
				useAABC = aabc
			}
			if d, ok := params["codeTableData"].([]byte); ok {
				customData = d
			}
		}
		if customData != nil {
			return NewWubiEncoder("custom", customData, useAABC)
		}
		return NewWubiEncoder(schema, nil, useAABC)
	default:
		return &NullEncoder{}
	}
}
