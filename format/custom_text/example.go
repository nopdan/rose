package custom_text

/*
新的 CustomText 格式使用示例

1. 基本字段类型：
   - FieldTypeWord: 词组
   - FieldTypePinyin: 拼音（支持分隔符、前后缀配置）
   - FieldTypeCode: 编码
   - FieldTypeFrequency: 词频
   - FieldTypeRank: 候选顺序
   - FieldTypeTab: 制表符 \t
   - FieldTypeSpace: 空格
   - FieldTypeLiteral: 自定义字面量

2. 示例：搜狗拼音格式（词组\tpin'yin\t词频）
   NewCustom(
       "sogou-pinyin",
       "搜狗拼音",
       model.FormatTypePinyin,
       util.NewEncoding("UTF-8"),
       []FieldConfig{
           {Type: FieldTypeWord},                          // 词组
           {Type: FieldTypeTab},                            // \t
           {Type: FieldTypePinyin, PinyinSeparator: "'"},  // pin'yin
           {Type: FieldTypeTab},                            // \t
           {Type: FieldTypeFrequency},                      // 词频
       },
       false,
       "#", // 注释前缀，以 # 开头的行会被跳过
   )

3. 示例：带前后缀的拼音格式（词组\t【pin'yin】\t词频）
   NewCustom(
       "custom-pinyin",
       "自定义拼音",
       model.FormatTypePinyin,
       util.NewEncoding("UTF-8"),
       []FieldConfig{
           {Type: FieldTypeWord},
           {Type: FieldTypeTab},
           {Type: FieldTypePinyin,
               PinyinSeparator: "'",
               PinyinPrefix: "【",
               PinyinSuffix: "】",
           },
           {Type: FieldTypeTab},
           {Type: FieldTypeFrequency},
       },
       false,
       "#",
   )

4. 示例：五笔格式（code\tword）
   NewCustom(
       "simple-wubi",
       "简单五笔",
       model.FormatTypeWubi,
       util.NewEncoding("UTF-8"),
       []FieldConfig{
           {Type: FieldTypeCode},
           {Type: FieldTypeTab},
           {Type: FieldTypeWord},
       },
       true, // 按编码排序
       "#",  // 支持注释
   )

5. 示例：自定义短语格式（abbr,rank=word）
   NewCustom(
       "user-phrase",
       "用户短语",
       model.FormatTypePinyin,
       util.NewEncoding("UTF-8"),
       []FieldConfig{
           {Type: FieldTypePinyin, PinyinSeparator: ""},  // 缩写无分隔符
           {Type: FieldTypeLiteral, Literal: ","},
           {Type: FieldTypeRank},
           {Type: FieldTypeLiteral, Literal: "="},
           {Type: FieldTypeWord},
       },
       false,
       "", // 不支持注释
   )

6. 示例：复杂自定义格式（word|pin'yin|freq|rank）
   NewCustom(
       "pipe-separated",
       "竖线分隔",
       model.FormatTypePinyin,
       util.NewEncoding("UTF-8"),
       []FieldConfig{
           {Type: FieldTypeWord},
           {Type: FieldTypeLiteral, Literal: "|"},
           {Type: FieldTypePinyin, PinyinSeparator: "'"},
           {Type: FieldTypeLiteral, Literal: "|"},
           {Type: FieldTypeFrequency},
           {Type: FieldTypeLiteral, Literal: "|"},
           {Type: FieldTypeRank},
       },
       false,
       "//", // 使用 // 作为注释前缀
   )

前端表单设计建议：
1. 提供字段类型下拉选择
2. 拼音字段额外提供：分隔符输入框、前缀输入框、后缀输入框
3. 字面量字段提供：内容输入框
4. 注释前缀输入框（可选，空表示不支持注释）
5. 支持拖拽排序字段
6. 提供快捷按钮：添加Tab、添加Space、添加逗号等
7. 实时预览生成的格式示例
*/
