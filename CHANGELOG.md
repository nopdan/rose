## [v1.3.1](https://github.com/nopdan/rose/compare/v1.3.0...v1.3.1) (2023-08-29)

### Ci

- 独立发布 data 文件夹

### Fix

- 匹配格式错误

## [v1.3.0](https://github.com/nopdan/rose/compare/v1.2.1...v1.3.0) (2023-08-29)

### Feat

- 卡饭五笔备份.dict
- 卡饭拼音备份.dict

### Perf

- 添加格式简写

## [v1.2.1](https://github.com/nopdan/rose/compare/v1.2.0...v1.2.1) (2023-08-29)

### Fix

- 自定义码表

### Perf

- 错误提示

## [v1.2.0](https://github.com/nopdan/rose/compare/v1.1.3...v1.2.0) (2023-08-29)

### Ci

- 分别构建 commit 和 release

### Feat

- 添加前端
- 自动编码五笔方案
- 支持多多生成器 v3 dmg 格式
- 支持多多生成器 v4 duodb 格式

### Perf

- 优化命令行

## [v1.1.3](https://github.com/nopdan/rose/compare/v1.1.2...v1.1.3) (2023-05-07)

### Fix

- 搜狗备份词库读取不完整

### Perf

- 提高注音准确度

## [v1.1.2](https://github.com/nopdan/rose/compare/v1.1.1...v1.1.2) (2023-04-19)

### Feat

- 支持微软五笔 lex 格式输出
- 支持微软拼音自学习词汇生成

### Fix

- 自学习过滤长词

### Perf

- 转换失败直接退出

## [v1.1.1](https://github.com/nopdan/rose/compare/v1.1.0...v1.1.1) (2023-04-18)

### Feat

- 自定义保存路径
- DetectFormat
- CodeTable.ToWubiTable

### Perf

- 添加一些简写

## [v1.1.0](https://github.com/nopdan/rose/compare/v1.0.1...v1.1.0) (2023-04-15)

### Style

- 抽象 Entry 接口

### Test

- 添加更多测试词库

## [v1.0.1](https://github.com/nopdan/rose/compare/v1.0...v1.0.1) (2023-04-15)

### Feat

- 添加交互式命令行
- 添加 `-v` 和 `-h` 命令

## [v1.0](https://github.com/nopdan/rose/compare/v0.7.3...v1.0) (2023-04-13)

### Feat

- 简易 cli
- 支持新的搜狗备份 bin 格式

### Perf

- sogou_scel black list
- 移除小胖输入法

### Refactor

- 删掉了大量功能，只保留词库转换
- 删掉了 GUI

## [v0.7.3](https://github.com/nopdan/rose/compare/v0.7.2...v0.7.3) (2023-03-14)

### Perf

- 支持小胖输入法二进制格式
- 极点排序

## [v0.7.2](https://github.com/nopdan/rose/compare/v0.7.1...v0.7.2) (2023-02-04)

### Feat

- 生成候选位置

## [v0.7.1](https://github.com/nopdan/rose/compare/v0.7...v0.7.1) (2023-02-02)

### Fix

- 双拼映射表 ve 韵母

## [v0.7](https://github.com/nopdan/rose/compare/v0.6...v0.7) (2022-09-17)

### Feat

- 微软自定义短语输出 & 五笔 lex 解析

## [v0.6](https://github.com/nopdan/rose/compare/v0.5...v0.6) (2022-09-16)

### Feat

- 全拼词库转为双拼定长码表
- 新的编码规则

## [v0.5](https://github.com/nopdan/rose/compare/v0.0.2...v0.5) (2022-09-12)

### Fix

- fcitx4 码表索引错误

### Perf

- 编码格式
- 搜狗 bin 解密拼音
- nothing
- 添加几个单字拼音
- 重写笛卡尔积

### Refactor

- 改接口，添加 checker
- 调整项目结构

## [v0.0.2](https://github.com/nopdan/rose/compare/v0.0.1...v0.0.2) (2022-06-01)

### Feat

- 出简不出全

### Refactor

- 规范命名，重组结构

## v0.0.1 (2022-05-31)

### Feat

- 简陋的 cli
- 写了个词库转换 GUI
- 独立 utils 文件夹
- 通过词表+单字表生成拼音
- 支持搜狗拼音备份 bin 格式
- 通用规则解析与生成
- 支持微软拼音自学习词汇.dat
- 支持微软拼音用户自定义短语.dat
- 支持最新的 scel 格式
- 支持含英文的.scel(.qcel)格式

### Perf

- 改进紫光 uwl 解析
- 不同系统使用不同换行符
- 使用反射
- 改了两个工具函数
- 使用一些位运算
- 优化异常处理
