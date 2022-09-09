## 目标

- [x] 码表格式转换（纯文本的好做，def 参考 asd2fque1 的 DictTool，其他参考深蓝）
- [x] 纯词生成拼音
- [x] 从全码表生成出简不出全码表（需要规则）
- [x] 从简码全码混和码表中提取单字全码
- [x] 根据单字全码表对词组进行编码（对多编码的字进行笛卡尔积，需要一种编码规则）
- [x] 根据单字全码表对词库进行错码校验
- [ ] 根据全拼词库生成双拼定长码表
- [ ] 词条过滤器（例：词长>9、码长>=5、词频<10）

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
