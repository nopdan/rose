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
