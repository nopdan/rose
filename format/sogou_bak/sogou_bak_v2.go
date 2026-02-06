package sogou_bak

import (
	"fmt"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

var utf16 = util.NewEncoding("UTF-16LE")

// keyItem 表示词典结构中的键配置项
// 对应Python中的KeyItem类
type keyItem struct {
	dictTypedef    uint16   // 词典类型定义
	dataTypeCount  uint16   // 数据类型数量
	dataTypes      []uint16 // 数据类型索引数组
	attributeIndex uint32   // 属性表索引
	keyDataIndex   uint32   // 键数据存储索引
	dataIndex      uint32   // 主数据存储索引
	reserved       uint32   // 保留字段
}

// attributeItem 表示属性配置项
// 对应Python中的AttributeItem类
type attributeItem struct {
	count          uint32 // 属性数量
	reserved1      uint32 // 保留字段1
	dataStoreIndex uint32 // 数据存储索引
	reserved2      uint32 // 保留字段2
}

// headerItem 表示数据段头，包含偏移量和大小信息
// 对应Python中的HeaderItem类
type headerItem struct {
	offset       uint32 // 数据段偏移量
	dataSize     uint32 // 总分配大小
	usedDataSize uint32 // 实际使用大小
}

// wordAttributeData 表示词汇元数据和加密信息
// 对应Python中的AttrWordData类
type wordAttributeData struct {
	textOffset    uint32 // 加密词汇文本的偏移量
	frequency     uint16 // 词频
	attributeFlag uint16 // 属性标志
	reserved1     uint32 // 保留字段1
	encryptionKey uint16 // 加密参数
	reserved2     uint32 // 保留字段2（总是为零）
}

// userDictionary 表示完整的词典结构
// 对应Python中的BaseDict类
type userDictionary struct {
	keyItems         []keyItem       // 键配置项
	attributeItems   []attributeItem // 属性配置项
	additionalInts   []uint32        // 附加整数值
	indexHeaders     []headerItem    // 索引段头
	attributeHeaders []headerItem    // 属性段头
	dataStoreHeaders []headerItem    // 数据存储段头

	// 初始化期间计算的派生值
	dataTypeSizes  []uint32 // 计算的数据类型大小
	attributeSizes []uint32 // 计算的属性大小
	baseHashSizes  []int    // 哈希表大小
	keyHashSizes   []int    // 键特定的哈希大小
	attributeFlag  bool     // 计算的属性标志
}

// 数据类型大小映射 - 对应Python中keyItem.datatype_size
var dataTypeSizeMap = []uint32{4, 1, 1, 2, 1, 2, 2, 4, 4, 8, 4, 4, 4, 0, 0, 0}

// 数据类型哈希大小映射 - 对应Python中BaseDict.datatype_hash_size
var dataTypeHashSizeMap = []int{0, 27, 414, 512, -1, -1, 512, 0}

// parseHeader 从读取器中读取头信息
func (h *headerItem) parseHeader(r *model.Reader) {
	h.offset = r.ReadUint32()
	h.dataSize = r.ReadUint32()
	h.usedDataSize = r.ReadUint32()
}

// parseWordInfo 从读取器中读取词汇属性数据
func (w *wordAttributeData) parseWordInfo(r *model.Reader) {
	w.textOffset = r.ReadUint32()
	w.frequency = r.ReadUint16()
	w.attributeFlag = r.ReadUint16()
	w.reserved1 = r.ReadUint32()
	w.encryptionKey = r.ReadUint16()
	w.reserved2 = r.ReadUint32() // 总是为零
	_ = r.ReadUint32()           // 跳过下一个偏移量
}

// newUserDictionary 创建新的用户词典实例
// 对应Python中BaseDict.__init__
func newUserDictionary() *userDictionary {
	userDict := &userDictionary{
		indexHeaders:     make([]headerItem, 0),
		attributeHeaders: make([]headerItem, 0),
		dataStoreHeaders: make([]headerItem, 0),
		dataTypeSizes:    make([]uint32, 0),
		baseHashSizes:    make([]int, 0),
		keyHashSizes:     make([]int, 10),
	}

	// 为第一个键初始化哈希大小（corev3兼容性）
	userDict.keyHashSizes[0] = 500
	return userDict
}

// importV2 解析搜狗备份词典v2格式
// 对应Python脚本中的主要解析逻辑
func (f *SogouBak) importV2(r *model.Reader) ([]*model.Entry, error) {
	f.Infof("开始解析搜狗备份词典v2格式\n")

	// 基本文件大小检查，对应Python中的size断言
	if r.Size() < 80 {
		return nil, fmt.Errorf("文件太小，不是有效的搜狗备份格式")
	}

	r.Seek(0, 0)

	// 解析文件头并获取配置信息
	keyCount, attrCount, aintCount := f.parseFileHeader(r)

	// 解析配置项
	keyItems, err := f.parseKeyItems(r, keyCount)
	if err != nil {
		return nil, err
	}
	attributeItems := f.parseAttributeItems(r, attrCount)
	additionalInts := f.parseAdditionalInts(r, aintCount)

	// 创建用户词典结构并设置配置
	userDict := f.setupUserDictionary(keyItems, attributeItems, additionalInts)

	// 解析数据段头信息
	f.parseDataHeaders(r, userDict)

	// 提取加密参数
	param2, param3 := f.extractEncryptionParams(r)

	// 初始化词典结构
	userDict.initialize()

	// 处理词汇条目
	entries := f.processWordEntries(r, userDict, param2, param3)

	f.Infof("成功解析了 %d 个条目\n", len(entries))
	return entries, nil
}

// parseFileHeader 解析文件头，返回各段的计数
func (f *SogouBak) parseFileHeader(r *model.Reader) (keyCount, attrCount, aintCount uint32) {
	f.Infof("解析文件头\n")

	// 读取Python版本中的文件大小和校验和
	fileSize := r.ReadUint32()
	f.Infof("文件大小: %d 字节\n", fileSize)

	checksum := r.ReadUint32()
	f.Infof("校验和: 0x%08x\n", checksum)

	keyCount = r.ReadUint32() // 键项数量
	f.Infof("键项数量: %d\n", keyCount)

	attrCount = r.ReadUint32() // 属性项数量
	f.Infof("属性项数量: %d\n", attrCount)

	aintCount = r.ReadUint32() // 附加整数数量
	f.Infof("附加整数数量: %d\n", aintCount)

	return keyCount, attrCount, aintCount
}

// parseKeyItems 解析键项配置
func (f *SogouBak) parseKeyItems(r *model.Reader, keyCount uint32) ([]keyItem, error) {
	f.Infof("解析 %d 个键项\n", keyCount)
	keyItems := make([]keyItem, 0, keyCount)

	for range keyCount {
		var keyItem keyItem
		keyItem.dictTypedef = r.ReadUint16()

		// 对应Python中的 assert key.dict_typedef < 100
		if keyItem.dictTypedef >= 100 {
			return nil, fmt.Errorf("键项字典类型定义异常: %d，应小于100", keyItem.dictTypedef)
		}

		keyItem.dataTypeCount = r.ReadUint16()
		keyItem.dataTypes = make([]uint16, keyItem.dataTypeCount)

		for j := range keyItem.dataTypeCount {
			keyItem.dataTypes[j] = r.ReadUint16()
			// 对应Python中的数据类型范围检查
			if keyItem.dataTypes[j] >= 16 {
				return nil, fmt.Errorf("键项[%d]数据类型[%d]异常: %d，应小于16", len(keyItems), j, keyItem.dataTypes[j])
			}
		}

		keyItem.attributeIndex = r.ReadUint32()
		keyItem.keyDataIndex = r.ReadUint32()
		keyItem.dataIndex = r.ReadUint32()
		keyItem.reserved = r.ReadUint32()

		f.Infof("键项[%d]: typedef=0x%x, dataTypes=%v, attrIdx=%d\n",
			len(keyItems), keyItem.dictTypedef, keyItem.dataTypes, keyItem.attributeIndex)
		keyItems = append(keyItems, keyItem)
	}

	return keyItems, nil
}

// parseAttributeItems 解析属性项
func (f *SogouBak) parseAttributeItems(r *model.Reader, attrCount uint32) []attributeItem {
	f.Infof("解析 %d 个属性项\n", attrCount)
	attributeItems := make([]attributeItem, 0, attrCount)

	for range attrCount {
		var attrItem attributeItem
		attrItem.count = r.ReadUint32()
		attrItem.reserved1 = r.ReadUint32()
		attrItem.dataStoreIndex = r.ReadUint32()
		attrItem.reserved2 = r.ReadUint32()

		f.Infof("属性项[%d]: count=%d, dataStoreIndex=%d\n",
			len(attributeItems), attrItem.count, attrItem.dataStoreIndex)
		attributeItems = append(attributeItems, attrItem)
	}

	return attributeItems
}

// parseAdditionalInts 解析附加整数
func (f *SogouBak) parseAdditionalInts(r *model.Reader, aintCount uint32) []uint32 {
	f.Infof("解析 %d 个附加整数\n", aintCount)
	additionalInts := make([]uint32, 0, aintCount)

	for range aintCount {
		aint := r.ReadUint32()
		additionalInts = append(additionalInts, aint)
	}

	return additionalInts
}

// setupUserDictionary 设置用户词典结构
func (f *SogouBak) setupUserDictionary(keyItems []keyItem, attributeItems []attributeItem, additionalInts []uint32) *userDictionary {
	userDict := newUserDictionary()
	userDict.keyItems = keyItems
	userDict.attributeItems = attributeItems
	userDict.additionalInts = additionalInts
	return userDict
}

// parseDataHeaders 解析数据段头信息
func (f *SogouBak) parseDataHeaders(r *model.Reader, userDict *userDictionary) {
	f.Infof("解析数据段头信息\n")

	// 读取Python版本中的版本和格式信息
	b2Ver := r.ReadUint32()
	f.Infof("数据版本: %d\n", b2Ver)

	b2Format := r.ReadUint32()
	f.Infof("数据格式: %d\n", b2Format)

	size2 := r.ReadUint32()
	f.Infof("数据段大小: %d\n", size2)

	indexHeaderCount := r.ReadUint32()
	f.Infof("索引头数量: %d\n", indexHeaderCount)

	attrHeaderCount := r.ReadUint32()
	f.Infof("属性头数量: %d\n", attrHeaderCount)

	dataStoreCount := r.ReadUint32()
	f.Infof("数据存储头数量: %d\n", dataStoreCount)

	// 解析索引头
	for range indexHeaderCount {
		var header headerItem
		header.parseHeader(r)
		userDict.indexHeaders = append(userDict.indexHeaders, header)
		f.Infof("索引头[%d]: offset=0x%x, dataSize=0x%x, usedSize=0x%x\n",
			len(userDict.indexHeaders)-1, header.offset, header.dataSize, header.usedDataSize)
	}

	// 解析属性头
	for range attrHeaderCount {
		var header headerItem
		header.parseHeader(r)
		userDict.attributeHeaders = append(userDict.attributeHeaders, header)
		f.Infof("属性头[%d]: offset=0x%x, dataSize=0x%x, usedSize=0x%x\n",
			len(userDict.attributeHeaders)-1, header.offset, header.dataSize, header.usedDataSize)
	}

	// 解析数据存储头
	for range dataStoreCount {
		var header headerItem
		header.parseHeader(r)
		userDict.dataStoreHeaders = append(userDict.dataStoreHeaders, header)
		f.Infof("数据存储头[%d]: offset=0x%x, dataSize=0x%x, usedSize=0x%x\n",
			len(userDict.dataStoreHeaders)-1, header.offset, header.dataSize, header.usedDataSize)
	}
}

// extractEncryptionParams 从文件末尾提取加密参数
func (f *SogouBak) extractEncryptionParams(r *model.Reader) (param2, param3 uint32) {
	// 对应Python版本中的p2和p3
	l := r.Len()
	r.Seek(r.Size()-0x14, 0)
	param2 = r.ReadUint32()
	param3 = r.ReadUint32()
	r.Seek(r.Size()-int64(l), 0) // 恢复到原始位置
	f.Infof("加密参数: p2=%d, p3=%d\n", param2, param3)
	return param2, param3
}

// processWordEntries 处理所有词汇条目
func (f *SogouBak) processWordEntries(r *model.Reader, userDict *userDictionary, encryptionParam2, encryptionParam3 uint32) []*model.Entry {
	// 获取当前偏移量用于数据访问
	baseDataOffset := r.Size() - int64(r.Len())
	f.Infof("基础数据偏移量: 0x%x\n", baseDataOffset)

	// 获取所有属性数据对进行处理
	attributeDataPairs := userDict.getAllAttributeData(r)
	f.Infof("找到 %d 个属性数据对\n", len(attributeDataPairs)/2)

	entries := make([]*model.Entry, 0, len(attributeDataPairs)/2)

	// 处理每个词汇条目
	for i := range len(attributeDataPairs) / 2 {
		indexDataOffset := attributeDataPairs[2*i]
		attributeDataOffset := attributeDataPairs[2*i+1]

		// 解密拼音代码
		r.Seek(int64(indexDataOffset), 0)
		pinyinPosition := r.ReadUint32()
		pinyinOffset := uint32(baseDataOffset) + userDict.dataStoreHeaders[userDict.keyItems[0].keyDataIndex].offset + pinyinPosition
		pinyinCodes := f.decryptPinyin(r, pinyinOffset)

		// 解析词汇信息
		var wordInfo wordAttributeData
		r.Seek(int64(attributeDataOffset), 0)
		wordInfo.parseWordInfo(r)

		// 解密词汇文本
		attributeIndex := userDict.keyItems[0].attributeIndex
		dataStoreIndex := userDict.attributeItems[attributeIndex].dataStoreIndex
		wordOffset := uint32(baseDataOffset) + userDict.dataStoreHeaders[dataStoreIndex].offset + wordInfo.textOffset
		word := decryptWordText(r, wordOffset, uint32(wordInfo.encryptionKey), encryptionParam2, encryptionParam3)

		f.Debugf("%s\t%v\t%d\n", word, pinyinCodes, wordInfo.frequency)
		entries = append(entries, model.NewEntry(word).WithMultiCode(pinyinCodes...).WithFrequency(int(wordInfo.frequency)))
	}

	return entries
}

// decryptPinyin 从指定偏移量解密拼音代码
// 对应Python中的DecryptPinyin函数
func (f *SogouBak) decryptPinyin(r *model.Reader, offset uint32) []string {
	r.Seek(int64(offset), 0)
	pinyinCount := r.ReadUint16() / 2
	pinyinCodes := make([]string, 0, pinyinCount)

	for range pinyinCount {
		pinyinIndex := r.ReadUint16()
		pinyinCodes = append(pinyinCodes, pyList[pinyinIndex])
	}

	return pinyinCodes
}

// decryptWordText 使用加密算法解密词汇文本
// 对应Python中的DecryptWordsEx函数
func decryptWordText(r *model.Reader, offset uint32, encryptionKey, param2, param3 uint32) string {
	// 计算加密密钥
	key1 := (encryptionKey + param2) << 2
	key2 := (encryptionKey + param3) << 2
	xorKey := (key1 + key2) & 0xFFFF

	r.Seek(int64(offset), 0)
	charCount := r.ReadUint16() / 2
	decryptedBytes := make([]byte, 0, charCount*2)

	for range charCount {
		shift := param2 % 8
		encryptedChar := uint32(r.ReadUint16())

		// 反向位旋转和异或解密
		decryptedChar := (encryptedChar<<(16-(shift%8)) | (encryptedChar >> shift)) & 0xFFFF
		decryptedChar ^= xorKey

		// 转换为小端字节
		decryptedBytes = append(decryptedBytes, byte(decryptedChar%0x100), byte(decryptedChar>>8))
	}

	// 解码UTF-16LE为字符串
	return utf16.Decode(decryptedBytes)
}

// initialize 计算词典结构的派生值
// 对应Python中BaseDict.init
func (ud *userDictionary) initialize() {
	ud.attributeSizes = make([]uint32, len(ud.attributeItems))

	for keyIndex, keyItem := range ud.keyItems {
		// 从词典类型计算基础大小
		size := (uint32(keyItem.dictTypedef) >> 2) & 4
		maskedType := int(keyItem.dictTypedef) & 0xFFFFFF8F

		// 确定哈希表大小
		if ud.keyHashSizes[keyIndex] > 0 {
			ud.baseHashSizes = append(ud.baseHashSizes, ud.keyHashSizes[keyIndex])
		} else {
			ud.baseHashSizes = append(ud.baseHashSizes, dataTypeHashSizeMap[maskedType])
		}

		// 计算数据类型大小
		attributeCount := ud.attributeItems[keyItem.attributeIndex].count
		nonAttributeCount := len(keyItem.dataTypes) - int(attributeCount)

		// 添加非属性数据类型大小
		for j := range nonAttributeCount {
			if j > 0 || maskedType != 4 {
				size += dataTypeSizeMap[keyItem.dataTypes[j]]
			}
		}

		// 添加条件大小调整
		if keyItem.dictTypedef&0x60 > 0 {
			size += 4
		}
		size += 4
		ud.dataTypeSizes = append(ud.dataTypeSizes, size)

		// 计算属性数据大小
		var attributeSize uint32
		for j := nonAttributeCount; j < len(keyItem.dataTypes); j++ {
			attributeSize += dataTypeSizeMap[keyItem.dataTypes[j]]
		}
		if (keyItem.dictTypedef & 0x40) == 0 {
			attributeSize += 4
		}
		ud.attributeSizes[keyItem.attributeIndex] = attributeSize

		// 根据保留字段设置属性标志
		if ud.attributeItems[keyItem.attributeIndex].reserved2 == 0 {
			ud.attributeFlag = true
		}
	}
}

// getAllAttributeData 提取所有属性数据对进行处理
// 对应Python中BaseDict.GetAllDataWithAttri
func (ud *userDictionary) getAllAttributeData(r *model.Reader) []uint32 {
	results := make([]uint32, 0, 0xff)
	keyIndex := 0
	keyItem := ud.keyItems[keyIndex]

	// 获取属性头信息
	attributeHeader := ud.attributeHeaders[keyItem.attributeIndex]
	var attributeCount uint32

	if attributeHeader.usedDataSize == 0 {
		attributeCount = attributeHeader.dataSize
	} else {
		attributeCount = attributeHeader.usedDataSize
	}
	hashStoreCount := ud.baseHashSizes[keyIndex]

	// 计算基础数据偏移量
	baseDataOffset := r.Size() - int64(r.Len())

	// 处理每个哈希存储
	for i := range hashStoreCount {
		// 读取哈希存储头
		r.Seek(baseDataOffset+int64(8*i), 0)
		hashStoreOffset := r.ReadUint32()
		hashStoreEntryCount := r.ReadUint32()

		// 处理哈希存储中的每个条目
		for j := range hashStoreEntryCount {
			// 计算属性偏移量
			attributeOffset := uint32(baseDataOffset) + ud.indexHeaders[keyIndex].offset +
				hashStoreOffset + ud.dataTypeSizes[keyIndex]*uint32(j)

			// 读取下一个属性偏移量
			nextOffsetPosition := attributeOffset + uint32(ud.dataTypeSizes[keyIndex]) - 4
			r.Seek(int64(nextOffsetPosition), 0)
			nextOffset := r.ReadUint32()

			// 处理属性链
			for range attributeCount {
				attributeDataOffset := uint32(baseDataOffset) +
					ud.attributeHeaders[ud.keyItems[keyIndex].attributeIndex].offset + nextOffset

				results = append(results, attributeOffset, attributeDataOffset)

				// 获取链中的下一个偏移量
				nextOffsetPosition = attributeDataOffset + ud.attributeSizes[keyItem.attributeIndex] - 4
				r.Seek(int64(nextOffsetPosition), 0)
				nextOffset = r.ReadUint32()

				// 检查链的结束
				if nextOffset == 0xFFFFFFFF {
					break
				}
			}
		}
	}

	return results
}
