package format

import (
	"fmt"

	"github.com/nopdan/rose/model"
)

// Registry 格式注册中心
type Registry struct {
	formats   map[string]model.Format
	aliases   map[string]string             // 别名映射
	typeIndex map[model.FormatType][]string // 按类型索引
	order     []string
}

// NewRegistry 创建新的注册中心
func NewRegistry() *Registry {
	return &Registry{
		formats:   make(map[string]model.Format),
		aliases:   make(map[string]string),
		typeIndex: make(map[model.FormatType][]string),
		order:     make([]string, 0),
	}
}

// GlobalRegistry 全局注册中心
var GlobalRegistry = NewRegistry()

// Register 注册格式
func (r *Registry) Register(format model.Format) {
	info := format.Info()
	if _, ok := r.formats[info.ID]; ok {
		fmt.Printf("Format %s already registered, skipping\n", info.ID)
		return
	}

	// 注册主格式
	r.formats[info.ID] = format
	r.order = append(r.order, info.ID)

	// 更新类型索引
	if r.typeIndex[info.Type] == nil {
		r.typeIndex[info.Type] = make([]string, 0)
	}
	r.typeIndex[info.Type] = append(r.typeIndex[info.Type], info.ID)

	// fmt.Printf("Registered format: %s (%s)\n", info.ID, info.Name)
}

// RegisterWithAliases 注册格式并添加别名
func (r *Registry) RegisterWithAliases(format model.Format, aliases ...string) {
	r.Register(format)

	info := format.Info()
	for _, alias := range aliases {
		if alias != "" && alias != info.ID {
			r.aliases[alias] = info.ID
		}
	}
}

// Get 获取格式处理器
func (r *Registry) Get(formatID string) (model.Format, bool) {
	// 先尝试直接查找
	if format, ok := r.formats[formatID]; ok {
		return format, true
	}

	// 再尝试别名查找
	if realID, ok := r.aliases[formatID]; ok {
		if format, ok := r.formats[realID]; ok {
			return format, true
		}
	}

	return nil, false
}

// List 获取所有格式信息
func (r *Registry) List() []*model.BaseFormat {
	infos := make([]*model.BaseFormat, 0, len(r.formats))
	for _, id := range r.order {
		format, ok := r.formats[id]
		if !ok {
			continue
		}
		info := format.Info()
		infos = append(infos, &model.BaseFormat{
			ID:          info.ID,
			Name:        info.Name,
			Type:        info.Type,
			Extension:   info.Extension,
			Description: info.Description,
		})
	}
	return infos
}

// ListByType 根据类型获取格式信息
func (r *Registry) ListByType(formatType model.FormatType) []*model.BaseFormat {
	infos := make([]*model.BaseFormat, 0)

	// 使用类型索引提高查找效率，保留注册顺序
	if ids, ok := r.typeIndex[formatType]; ok {
		for _, id := range ids {
			if format, ok := r.formats[id]; ok {
				info := format.Info()
				infos = append(infos, &model.BaseFormat{
					ID:          info.ID,
					Name:        info.Name,
					Type:        info.Type,
					Extension:   info.Extension,
					Description: info.Description,
				})
			}
		}
	}
	return infos
}

// RegisterFormat 向全局注册中心注册格式（向后兼容）
func RegisterFormat(format model.Format) {
	GlobalRegistry.Register(format)
}

// RegisterFormatWithAliases 向全局注册中心注册格式并添加别名
func RegisterFormatWithAliases(format model.Format, aliases ...string) {
	GlobalRegistry.RegisterWithAliases(format, aliases...)
}
