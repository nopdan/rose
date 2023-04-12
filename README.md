# 蔷薇词库转换

[![GitHub Repo stars](https://img.shields.io/github/stars/flowerime/rose)](https://github.com/flowerime/rose/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/flowerime/rose)](https://github.com/flowerime/rose/network/members)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/flowerime/rose)](https://github.com/flowerime/rose/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/flowerime/rose/build.yml)](https://github.com/flowerime/rose/actions/workflows/build.yml)
![GitHub repo size](https://img.shields.io/github/repo-size/flowerime/rose)
![GitHub](https://img.shields.io/github/license/flowerime/rose)

关于词库格式的详细解析可以到我的[博客](https://nopdan.com/)查看。

## 拼音词库格式

**纯文本：**

| 描述     | 代号      | 词条格式                     | 编码格式 |
| :------- | :-------- | :--------------------------- | -------- |
| 搜狗拼音 | sogou     | 拼音('分隔)` `词             |          |
| qq 拼音  | qq        | 拼音('分隔)` `词` `词频      |          |
| 百度拼音 | baidu     | 词`\t`拼音('分隔)`\t`词频    |          |
| 谷歌拼音 | google    | 词`\t`词频`\t`拼音(空格分隔) |          |
| 拼音加加 | pyjj      | 字音字音字音...              |          |
| 纯汉字   | word_only | 一词一行                     |          |

> 纯汉字词库会[自动注音](./pkg/zhuyin/zhuyin.go)，所以你可以当成注音工具使用。

**二进制：**

> **加粗**项支持输出

| 描述                         | 代号        | 格式             |
| :--------------------------- | :---------- | :--------------- |
| 搜狗细胞词库、qq6.0 以上词库 | sogou_scel  | `.scel`\|`.qcel` |
| 搜狗拼音备份词库             | sogou_bin   | `.bin`           |
| qq6.0 以下词库               | qq_qpyd     | `.qpyd`          |
| 百度分类词库                 | baidu_bdict | `.bdict`\|`.bcd` |
| 紫光（华宇）                 | ziguang_uwl | `.uwl`           |
| **微软用户自定义短语**       | mspy_dat    | `.dat`           |
| 微软拼音自学习词汇           | mspy_udl    | `.dat`           |

## 字词码表格式

**纯文本：**

| 描述 | 代号     | 词条格式               | 编码格式 |
| :--- | :------- | :--------------------- | -------- |
| 多多 | duoduo   | 词`\t`编码             |          |
| 冰凌 | bingling | 编码`\t`词             | UTF-16LE |
| 极点 | jidian   | 编码`\t`词 1` `词 2... |          |

**二进制：**

> **加粗**项支持输出

| 描述                   | 代号      | 格式   |
| :--------------------- | :-------- | :----- |
| **百度手机自定义方案** | baidu_def | `.def` |
| **微软用户自定义短语** | msudp_dat | `.dat` |
| 极点                   | jidian_mb | `.mb`  |
| fcitx4                 | fcitx4_mb | `.mb`  |
