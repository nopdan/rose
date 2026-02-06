# 测试数据文件

本目录包含项目测试所需的样例数据文件。

## 文件说明

### 拼音输入法格式
- `test.scel` - 搜狗细胞词库测试文件
- `hotcell_dict.scel` - 搜狗热词词库
- `sogou.scel` - 搜狗词库样例
- `成语_官方.scel` - 搜狗成语词库
- `网络流行新词.scel` - 搜狗网络流行词
- `qq.qcel` - QQ拼音细胞词库
- `计算机硬件词汇.bdict` - 百度分类词库
- `baidu.bcd` - 百度词库样例
- `ChsPinyinUDP.lex` - 微软拼音自学习词汇
- `拼音词库备份2023-08-19.dict` - 卡饭拼音备份词库
- `拼音词库备份2023-08-29.dict` - 卡饭拼音备份词库（新版）
- `music.uwl` - 紫光华宇拼音词库
- `sogou_bak_v2.bin` - 搜狗备份词库 v2
- `sogou_bak_v3.bin` - 搜狗备份词库 v3
- `jiajia.txt` - 拼音加加格式样例
- `rime_sample.dict.yaml` - Rime 词典格式样例
- `ChsPinyinUDL.dat` - 微软拼音用户词典

### 五笔输入法格式
- `ChsWubi.lex` - 微软五笔词典
- `五笔词库备份2023-08-29.dict` - 卡饭五笔备份词库
- `sample.def` - 百度手机自定义方案
- `98wb_ci.mb` - Fcitx4 码表
- `jidian.mb` - 极点码表
- `UserDefinedPhrase.dat` - 微软用户自定义短语

### QQ拼音格式
- `qq.qpyd` - QQ拼音分类词库
- `网络用语.qpyd` - QQ拼音网络用语词库

### 其他
- `words_*.txt` - 纯词组格式临时测试文件

## 使用说明

这些文件被各个格式的测试文件引用，路径统一为相对于项目根目录的 `testdata/` 目录。

测试时会自动从该目录读取相应的样例文件。
