<div align="center">

# 蔷薇词库转换

[![GitHub Repo stars](https://img.shields.io/github/stars/nopdan/rose)](https://github.com/nopdan/rose/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/nopdan/rose)](https://github.com/nopdan/rose/network/members)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/nopdan/rose)](https://github.com/nopdan/rose/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/nopdan/rose/commit.yml)](https://github.com/nopdan/rose/actions/workflows/commit.yml)
![GitHub repo size](https://img.shields.io/github/repo-size/nopdan/rose)
![GitHub](https://img.shields.io/github/license/nopdan/rose)

关于词库格式的详细解析可以到我的[博客](https://nopdan.com/)查看。

</div>

## 设计目标

支持的：

1. 转换各个输入法的私有格式备份词库，方便用户迁移输入法。
2. 转换大厂输入法的词库（如搜狗细胞词库，百度分类词库），导入小厂输入法使用。
3. 其他词库转五笔——需要选择不同五笔方案或自定义。
4. 其他词库转拼音——需要实现自动注音。
5. [TODO]过滤。根据词长，词频，是否含英文等条件过滤。

不支持的：

1. 英文词典，简繁转换，文件分割，自动爬取词频等。
2. 自动添加 Rime，小小，极点等文件头（意思就是你要手动添加）。
3. 其他词库转五笔，只支持四码定长的形码方案。不支持更加高级的选项，例如根据拼音转换为四码定长的双拼词库、二笔词库，类似键道 6 的六码方案，红辣椒五笔的不定长形码，出简不出全，码表合并等。若有此类需求可以去看我的另一个项目 [lilac](https://github.com/nopdan/lilac)。
4. 小胖输入法（作者不支持，不想与其斗智斗勇）。

词库形式：

1. 拼音词库。词组，分隔符分隔的拼音，可能有词频。
2. 五笔码表。词组，编码，可能有候选位置。
3. 用户自定义短语。词组，编码，可能有候选位置。
4. 纯词组。

优先支持：windows 平台，拼音词库，备份词库。

## 使用

```sh
Root Command:
    Usage: rose [输入文件] [输入格式]:[输出格式] [保存文件名]
    Example: rose sogou.scel scel:rime rime.dict.yaml

Sub Commands:
      list      列出所有支持的格式
      server    启动服务  -p:[port] 指定端口(默认7800)
  -h, help      帮助
  -v, version   版本
```

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

## 字词码表格式

| 词库               | 代号      | 简写 | 格式                   | 备注       |
| ------------------ | --------- | ---- | ---------------------- | ---------- |
| 百度手机自定义方案 | baidu_def | def  | `.def`                 |            |
| 微软用户自定义短语 | msudp_dat | udp  | `.dat`                 |            |
| 微软五笔           | mswb_lex  | lex  | `.lex`                 |            |
| 极点               | jidian_mb | jdmb | `.mb`                  | 不支持输出 |
| fcitx4             | fcitx4_mb | f4mb | `.mb`                  | 不支持输出 |
| 多多 v3            | dddmg     | dmg  | `.dmg`                 | 不支持输出 |
| 多多 v4            | duodb     |      | `.duodb`               | 不支持输出 |
| 多多               | duoduo    | dd   | 词`\t`编码             |            |
| 冰凌               | bingling  | bl   | 编码`\t`词             | UTF-16LE   |
| 极点               | jidian    | jd   | 编码`\t`词 1` `词 2... |            |

## 编译

```powershell
cd build
.\build.ps1
```
