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

- 转换各个输入法的私有格式备份词库，方便用户迁移输入法。
- 转换大厂输入法的词库（如搜狗细胞词库，百度分类词库），导入小厂输入法使用。
- 其他词库转五笔——需要选择不同五笔方案或自定义。
- 其他词库转拼音——需要实现自动注音。
- [TODO]过滤。根据词长，词频，是否含英文等条件过滤。

不支持的：

- 英文词典，简繁转换，文件分割，自动爬取词频等。
- 自动添加 Rime，小小，极点等文件头（意思就是你要手动添加）。
- 其他词库转五笔，只支持四码定长的形码方案。不支持更加高级的选项，例如根据拼音转换为四码定长的双拼词库、二笔词库，类似键道 6 的六码方案，红辣椒五笔的不定长形码，出简不出全，码表合并等。若有此类需求可以去看我的另一个项目 [lilac](https://github.com/nopdan/lilac)。
- 小胖输入法（作者不支持，不想与其斗智斗勇）。

词库形式：

1. 拼音词库。词组，分隔符分隔的拼音，可能有词频。
2. 五笔码表。词组，编码，可能有候选位置。
3. 用户自定义短语。词组，编码，可能有候选位置。
4. 纯词组。

优先支持：windows 平台，拼音词库，备份词库。

## 使用

**下载 data.zip 并解压，确保 data 文件夹与 rose 程序位于同一目录下。**

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

## 支持格式

| ID                  | 格式                   | 可导出 |
| ------------------- | ---------------------- | ------ |
| baidu,bd            | 百度拼音               | 是     |
| baidu_bcd,bcd       | 百度手机分类词库.bcd   |        |
| baidu_bdict,bdict   | 百度分类词库.bdict     |        |
| baidu_def,def       | 百度手机自定义方案.def | 是     |
| bingling,bl         | 冰凌                   | 是     |
| duoduo,dd           | 多多                   | 是     |
| duoduo_dmg,dmg      | 多多 v3.dmg            |        |
| duoduo_duodb,duodb  | 多多 v4.duodb          |        |
| google,gg           | 谷歌拼音               | 是     |
| jidian,jd           | 极点码表               | 是     |
| jidian_mb,jdmb,mb   | 极点码表.mb            |        |
| kfpybak,dict        | 卡饭拼音备份.dict      |        |
| kfwbbak             | 卡饭五笔备份.dict      |        |
| ms_udp,udp,dat      | 微软用户自定义短语.dat | 是     |
| mspy_udl,udl        | 微软拼音自学习词汇.dat | 是     |
| mswb_lex,lex        | 微软五笔.lex           | 是     |
| pyjj,jj             | 拼音加加               | 是     |
| qq                  | QQ 拼音                | 是     |
| qq_qcel,qcel        | QQ 拼音 v6 以上.qcel   |        |
| qq_qpyd,qpyd        | QQ 拼音 v6 以下.qpyd   |        |
| rime                | Rime 拼音              | 是     |
| sogou,sg            | 搜狗拼音               | 是     |
| sogou_bak,sgbak,bin | 搜狗拼音备份.bin       |        |
| sogou_scel,scel     | 搜狗细胞词库.scel      |        |
| words               | 纯词组                 | 是     |
| ziguang_uwl,uwl     | 紫光华宇拼音.uwl       |        |

## 编译

```powershell
git clone https://github.com/nopdan/rose.git
cd rose
git submodule update --init

cd build
.\build.ps1
```

## 架构与实现说明

- 架构说明：见 [docs/architecture.md](docs/architecture.md)
- 格式实现情况：见 [docs/format-status.md](docs/format-status.md)
