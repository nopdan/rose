# 词库处理工具

[![GitHub Repo stars](https://img.shields.io/github/stars/imetool/dtool)](https://github.com/imetool/dtool/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/imetool/dtool)](https://github.com/imetool/dtool/network/members)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/imetool/dtool)](https://github.com/imetool/dtool/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/imetool/dtool/build.yml)](https://github.com/imetool/dtool/actions/workflows/build.yml)
![GitHub repo size](https://img.shields.io/github/repo-size/imetool/dtool)
![GitHub](https://img.shields.io/github/license/imetool/dtool)

词库处理工具，词库编码，词库格式转换，词库校验，出简不出全。

关于词库格式的详细解析可以到我的[博客](https://noif.cc/)查看。

## [拼音词库转换](./pkg/pinyin/)

**纯文本：**

| 描述     | 代号      | 词条格式                     | 编码格式 |
| :------- | :-------- | :--------------------------- | -------- |
| 搜狗拼音 | sogou     | 拼音('分隔)` `词             |          |
| qq 拼音  | qq        | 拼音('分隔)` `词` `词频      |          |
| 百度拼音 | baidu     | 词`\t`拼音('分隔)`\t`词频    |          |
| 谷歌拼音 | google    | 词`\t`词频`\t`拼音(空格分隔) |          |
| 拼音加加 | pyjj      | 字音字音字音...              |          |
| 纯汉字   | word_only | 一词一行                     |          |

> 纯汉字词库会[自动注音](./pkg/encoder/pinyin.go)，所以你可以当成注音工具使用。

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

## [字词码表转换](./pkg/table/)

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

## [转双拼词库](./pkg/double/double_pinyin.go)

全拼词库转为双拼四码定长码表，需要一个双拼映射表，格式为手心输入法，具体可以看[示例](assets/双拼映射表/)。

## [词库编码和校验](./pkg/checker/checker.go)

> 仅支持多多格式，其他格式需先行转换

编码选择纯词词库，校验选择带编码的多多格式词库。

形如 `2=AaAbBaBb,3=AaAbBaCa,0=AaBaCaZa` ，等号前数字表示词长，0 表示未指定的词长。

也可以简写为 `2=AABB,3=AABC,0=ABCZ`

对于整句，`ab...`(必须以...结尾) 表示取每个字编码的前两码

## [全码转简码码表](./pkg/encoder/shortener.go)

> 仅支持多多格式，其他格式需先行转换

出简不出全规则：逗号，冒号分隔，默认 1，n 无限

```yaml
例子: '1:0, 2:3, 3:2, 6:n'
#  无 1 简，2 码 3 重，3 码 2 重，4 码 1 重，5 码 1 重，6 码无限重
```
