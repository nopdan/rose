<div align="center">

# 蔷薇词库转换

[![GitHub Repo stars](https://img.shields.io/github/stars/nopdan/rose)](https://github.com/nopdan/rose/stargazers)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/nopdan/rose)](https://github.com/nopdan/rose/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/nopdan/rose/commit.yml)](https://github.com/nopdan/rose/actions/workflows/commit.yml)
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
| `msudp` | 微软用户自定义短语 | .dat | ✓ |
| `mswb_lex` | 微软五笔词典 | .lex | ✓ |
| `jidian_mb` | 极点码表.mb | .mb | |
| `fcitx4_mb` | Fcitx4 码表 | .mb | |
| `kafan_wubi_bak` | 卡饭五笔备份 | .dict | |
| `jidian` | 极点码表 | .txt | ✓ |
| `duoduo` | 多多生成器 | .txt | ✓ |
| `baidu_wubi` | 百度五笔 | .txt | ✓ |
| `user_phrase` | 用户自定义短语 | .txt | ✓ |

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

## License

[MIT](LICENSE)
