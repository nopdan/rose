## 逻辑

转换：

1. 选择源码表
2. 选择源码表格式
3. 根据源码表格式读取为 WcTable, WpfDict。
4. 执行过滤器
5. 选择是否转换目标方案（纯词、全拼、双拼、形码）
6. 选择目标方案依赖项（双拼依赖一个映射表、形码依赖一个单字全码表）
7. 转换方案，全拼 WpfDict，双拼 WcTable，形码 WcTable
8. 目标格式，dict 反射为 WpfDict 和 WcTable

字词方案的工具：

1. 选择源码表
2. 选择源码表格式（WcTable, CwsTable）
3. 选择要进行的操作（出简、检验）
4. （检验需要单字全码表）
5. just do it

## 规范

词：Word string
单字：Char rune
拼音（多个音节）：Pinyin []string
编码|音节：code string
多个编码|读音：codes []string
词频：Freq

词长：wordLen
词占用字节数：wordSize
拼音长（音节数）：pyLen
拼音占用字节数：pySize
编码长度：codeLen (=codeSize)

文件中的字节位置，偏移量 pos
在候选词条中的顺序，码表中的第几行 order

字词->码表 Table
拼音->词库 Dict

一行，词条 entry
