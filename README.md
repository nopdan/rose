<div align="center">

# 蔷薇词库转换

[![GitHub Repo stars](https://img.shields.io/github/stars/nopdan/rose)](https://github.com/nopdan/rose/stargazers)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/nopdan/rose)](https://github.com/nopdan/rose/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/nopdan/rose/commit.yml)](https://github.com/nopdan/rose/actions/workflows/commit.yml)
![Downloads](https://badgen.net/github/assets-dl/nopdan/rose)
![GitHub](https://img.shields.io/github/license/nopdan/rose)

多种输入法词库格式互相转换，支持拼音、五笔、纯词组。

</div>

## 功能

- 转换各大输入法的私有格式备份词库，方便用户迁移输入法
- 转换搜狗细胞词库、百度分类词库等，导入其他输入法使用
- 拼音词库自动注音
- 五笔词库编码转换，支持 86/98/新世纪方案及自定义码表
- 按词长、词频、是否含英文/数字等条件过滤
- 自定义纯文本格式（灵活配置字段排列、分隔符、编码等）
- 内置 Web 界面，上传文件后自动识别格式、一键转换下载

## 使用

### Web 界面

```sh
rose serve        # 启动 Web 服务，默认端口 7800
rose serve 8080   # 指定端口
```

浏览器自动打开 `http://localhost:7800`，按页面引导上传词库文件、选择输入/输出格式即可完成转换。

### 命令行

```sh
rose <输入文件> <输入格式> <输出格式> [输出文件]
```

示例：

```sh
rose sogou.scel sogou_scel rime output.dict.yaml
rose backup.bin sogou_bak baidu baidu.txt
```

其他命令：

```sh
rose -l   # 列出所有支持的格式
rose -v   # 版本信息
rose -h   # 帮助
```

## 支持格式

共 25 种格式，全部支持导入，17 种支持导出。

### 拼音

| ID | 名称 | 扩展名 | 导出 |
|---|---|---|:---:|
| `sogou_scel` | 搜狗细胞词库 | .scel | ✓ |
| `qq_qcel` | QQ 拼音 v6+.qcel | .qcel | ✓ |
| `sogou_bak` | 搜狗拼音备份 | .bin | |
| `baidu_bdict` | 百度分类词库 | .bdict | |
| `baidu_bcd` | 百度手机分类词库 | .bcd | |
| `baidu_bak` | 百度拼音备份 | .bin | |
| `qq_qpyd` | QQ 拼音分类词库 | .qpyd | |
| `mspy_udl` | 微软拼音自学习词汇 | .dat | ✓ |
| `ziguang_uwl` | 紫光华宇拼音词库 | .uwl | |
| `kafan_pinyin_bak` | 卡饭拼音备份 | .dict | |
| `jiajia` | 拼音加加 | .txt | ✓ |
| `rime` | Rime 输入法 | .dict.yaml | ✓ |
| `sogou_pinyin` | 搜狗拼音 | .txt | ✓ |
| `baidu` | 百度拼音 | .txt | ✓ |
| `qq` | QQ 拼音 | .txt | ✓ |

### 五笔

| ID | 名称 | 扩展名 | 导出 |
|---|---|---|:---:|
| `baidu_def` | 百度手机自定义方案 | .def | ✓ |
| `mswb_lex` | 微软五笔词典 | .lex | ✓ |
| `jidian_mb` | 极点码表.mb | .mb | |
| `fcitx4_mb` | Fcitx4 码表 | .mb | |
| `kafan_wubi_bak` | 卡饭五笔备份 | .dict | |
| `jidian` | 极点码表 | .txt | ✓ |
| `duoduo` | 多多生成器 | .txt | ✓ |
| `baidu_shouji` | 百度手机 | .txt | ✓ |

### 自定义短语

| ID | 名称 | 扩展名 | 导出 |
|---|---|---|:---:|
| `user_phrase` | 用户自定义短语 | .txt | ✓ |
| `msudp` | 微软用户自定义短语 | .dat | ✓ |

### 纯词组

| ID | 名称 | 扩展名 | 导出 |
|---|---|---|:---:|
| `words` | 纯词组 | .txt | ✓ |

## 编译

需要 Go 1.24+，前端需要 bun.js。

```sh
git clone https://github.com/nopdan/rose.git
cd rose

# 构建前端（可选，已内嵌）
cd frontend && bun install && bun run build && cd ..

# 编译
go build .
```

## 项目结构

```
converter/    核心转换流程（导入 → 编码 → 过滤 → 导出）
encoder/      编码器（拼音注音、五笔编码）
filter/       过滤器（词长、词频、正则等）
format/       格式注册表及所有格式实现
model/        数据模型（Entry、Format、Importer/Exporter 接口）
server/       HTTP API 服务
frontend/     Vue 3 + Naive UI Web 界面
```

## 贡献新格式

添加一种新的词库格式只需三步：

### 1. 创建格式包

在 `format/` 下新建目录，例如 `format/my_format/`，创建主文件：

```go
package my_format

import (
    "io"
    "github.com/nopdan/rose/model"
)

type MyFormat struct {
    model.BaseFormat
}

func New() *MyFormat {
    return &MyFormat{
        BaseFormat: model.BaseFormat{
            ID:        "my_format",       // 唯一标识，用于命令行和 API
            Name:      "我的格式",         // 显示名称
            Type:      model.FormatTypePinyin, // FormatTypePinyin / FormatTypeWubi / FormatTypeWords
            Extension: ".txt",            // 文件扩展名（含点号）
        },
    }
}
```

### 2. 实现 Importer / Exporter 接口

按需实现导入和导出，至少实现其中一个：

```go
// 导入：从 Source 读取数据，返回 Entry 列表
func (f *MyFormat) Import(src model.Source) ([]*model.Entry, error) {
    // OpenTextReader 自动识别编码格式转为 utf8 reader
    // textReader, _, closeFn, err := model.OpenTextReader(src)
    // 获取 Reader
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}
    // 解析 data，构建 Entry 列表
    // 每个 Entry 包含 Word（词组）、Code（编码）、Frequency（词频）等字段
    var entries []*model.Entry
    // ...解析逻辑
    return entries, nil
}

// 导出：将 Entry 列表写入 Writer
func (f *MyFormat) Export(entries []*model.Entry, w io.Writer) error {
    // 将 entries 按目标格式写入 w
    return nil
}
```

### 3. 注册格式

在 `format/init.go` 的 `init()` 中添加注册：

```go
import "github.com/nopdan/rose/format/my_format"

// 在对应分组中添加
RegisterFormat(my_format.New())
```

注册后，该格式自动在命令行（`rose -l`）和 Web 界面中可用。

### 建议

- 参考 `format/baidu_bdict/`（二进制导入）或 `format/rime/`（纯文本导入导出）作为模板
- 如果是纯文本且格式简单（词条 + 分隔符 + 编码/拼音），可以直接用 `custom_text.NewCustom()` 创建，无需新建包
- 在 `testdata/` 下放入测试用的样本文件，编写 `*_test.go` 验证导入导出正确性

## License

[GPL v3](LICENSE)
