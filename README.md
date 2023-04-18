<div align="center">

<img src="logo.png"  width="150" height="150"> </img>

### 蔷薇词库转换

[![GitHub Repo stars](https://img.shields.io/github/stars/flowerime/rose)](https://github.com/flowerime/rose/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/flowerime/rose)](https://github.com/flowerime/rose/network/members)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/flowerime/rose)](https://github.com/flowerime/rose/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/flowerime/rose/build.yml)](https://github.com/flowerime/rose/actions/workflows/build.yml)
![GitHub repo size](https://img.shields.io/github/repo-size/flowerime/rose)
![GitHub](https://img.shields.io/github/license/flowerime/rose)

关于词库格式的详细解析可以到我的[博客](https://nopdan.com/)查看。

</div>

## 使用

`.\rose.exe input input_format:output_format output`

> 例：`.\rose.exe .\sogou.scel scel:rime rime.dict.yaml`

## 拼音词库格式

| 词库                         | 代号        | 简写  | 格式                         | 备注       |
| ---------------------------- | ----------- | ----- | ---------------------------- | ---------- |
| 搜狗细胞词库、qq6.0 以上词库 | sogou_scel  | scel  | `.scel`\|`.qcel`             | 不支持输出 |
| 搜狗拼音备份词库             | sogou_bin   | sgbin | `.bin`                       | 不支持输出 |
| qq6.0 以下词库               | qq_qpyd     | qpyd  | `.qpyd`                      | 不支持输出 |
| 百度分类词库                 | baidu_bdict | bdict | `.bdict`\|`.bcd`             | 不支持输出 |
| 紫光（华宇）                 | ziguang_uwl | uwl   | `.uwl`                       | 不支持输出 |
| 微软用户自定义短语           | mspy_dat    | udp   | `.dat`                       |            |
| 微软拼音自学习词汇           | mspy_udl    | udl   | `.dat`                       |            |
| 搜狗拼音                     | sogou       | sg    | 拼音('分隔)` `词             |            |
| qq 拼音                      | qq          |       | 拼音('分隔)` `词` `词频      |            |
| 百度拼音                     | baidu       | bd    | 词`\t`拼音('分隔)`\t`词频    |            |
| 谷歌拼音                     | google      | gg    | 词`\t`词频`\t`拼音(空格分隔) |            |
| rime                         | rime        |       | 词`\t`拼音(空格分隔)`\t`词频 |            |
| 拼音加加                     | jiajia      | jj    | 字音字音字音...              |            |
| 纯汉字                       | word_only   | w     | 一词一行                     |            |

> 纯汉字词库会[自动注音](./pkg/zhuyin/zhuyin.go)，所以可以当做注音工具使用。

## 字词码表格式

| 词库               | 代号      | 简写 | 格式                   | 备注       |
| ------------------ | --------- | ---- | ---------------------- | ---------- |
| 百度手机自定义方案 | baidu_def | def  | `.def`                 |            |
| 微软用户自定义短语 | msudp_dat | udp  | `.dat`                 |            |
| 微软五笔           | mswb_lex  | lex  | `.lex`                 | 不支持输出 |
| 极点               | jidian_mb | jdmb | `.mb`                  | 不支持输出 |
| fcitx4             | fcitx4_mb | f4mb | `.mb`                  | 不支持输出 |
| 多多               | duoduo    | dd   | 词`\t`编码             |            |
| 冰凌               | bingling  | bl   | 编码`\t`词             | UTF-16LE   |
| 极点               | jidian    | jd   | 编码`\t`词 1` `词 2... |            |
