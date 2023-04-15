## 规范

|                          |   变量名   | 数据类型(没标就是 int) |
| ------------------------ | :--------: | :--------------------: |
| 文件中的字节位置，偏移量 |   offset   |                        |
| 字词码表                 |   table    |           #            |
| 拼音词库                 |    dict    |           #            |
| 词条                     |   entry    |           #            |
| 词条数                   | entryCount |                        |
| 单字                     |    char    |          rune          |
| 词                       |    word    |         string         |
| 词长                     |  wordLen   |                        |
| 词占用字节数             |  wordSize  |                        |
| 编码                     |    code    |         string         |
| 编码长度                 |  codeLen   |                        |
| 音节                     |   yinjie   |         string         |
| 拼音（多个音节）         |   pinyin   |        []string        |
| 拼音长（音节数）         |   pyLen    |                        |
| 拼音占用字节数           |   pySize   |                        |
| 在候选中的位置           |    pos     |                        |
| 词频                     |    freq    |                        |
