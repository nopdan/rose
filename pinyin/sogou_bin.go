package pinyin

import (
	"bytes"
	"io"
	"io/ioutil"

	. "github.com/cxcn/dtool/utils"
)

type key struct {
	dictType    int // 2 字节
	dataTypeLen int // 2 字节
	dataType    []int
	attrIdx     int // 4 字节
	keyDataIdx  int // 4 字节
	dataIdx     int // 4 字节
	v6          int // 4 字节
}

// 全是 4 字节
type attr struct {
	count  int
	a2     int
	dataId int
	b2     int
}

// 全是 4 字节
type header struct {
	offset       int // 偏移量
	dataSize     int
	usedDataSize int
}

func (h *header) parse(r *bytes.Reader) {
	h.offset = ReadInt(r, 4)
	h.dataSize = ReadInt(r, 4)
	h.usedDataSize = ReadInt(r, 4)
}

func ParseSogouBin(rd io.Reader) []PyEntry {
	ret := make([]PyEntry, 0, 0xff)
	data, _ := ioutil.ReadAll(rd)
	r := bytes.NewReader(data)
	// var tmp []byte

	// fileChksum := ReadInt(r, 4)
	// size1 := ReadInt(r, 4)
	r.Seek(8, 1)
	keyLen := ReadInt(r, 4)
	attrLen := ReadInt(r, 4)
	aintLen := ReadInt(r, 4)
	// fmt.Println(fileChksum, size1, keyLen, attrLen, aintLen)

	keys := make([]key, 0, 1)
	for i := 0; i < keyLen; i++ {
		var k key
		k.dictType = ReadInt(r, 2)
		k.dataTypeLen = ReadInt(r, 2)
		k.dataType = make([]int, 0, 1)
		for j := 0; j < k.dataTypeLen; j++ {
			dataType := ReadInt(r, 2)
			k.dataType = append(k.dataType, dataType)
		}

		k.attrIdx = ReadInt(r, 4)
		k.keyDataIdx = ReadInt(r, 4)
		k.dataIdx = ReadInt(r, 4)
		k.v6 = ReadInt(r, 4)

		keys = append(keys, k)
	}

	attrs := make([]attr, 0, 1)
	for i := 0; i < attrLen; i++ {
		var a attr
		a.count = ReadInt(r, 4)
		a.a2 = ReadInt(r, 4)
		a.dataId = ReadInt(r, 4)
		a.b2 = ReadInt(r, 4)
		attrs = append(attrs, a)
	}

	aints := make([]int, 0, 1)
	for i := 0; i < aintLen; i++ {
		aint := ReadInt(r, 4)
		aints = append(aints, aint)
	}

	// fmt.Printf("keys %+v\n", keys)
	// fmt.Printf("attrs %+v\n", attrs)
	// fmt.Printf("aints %+v\n", aints)

	ud := newUsrDict()
	ud.keys = keys
	ud.attrs = attrs
	ud.aints = aints

	// b2Ver := ReadInt(r, 4)
	// b2Format := ReadInt(r, 4)
	// size2 := ReadInt(r, 4)
	r.Seek(12, 1)
	// fmt.Println(b2Ver, b2Format, size2)

	hiLen := ReadInt(r, 4)
	haLen := ReadInt(r, 4)
	hsLen := ReadInt(r, 4)
	// fmt.Println(hiLen, haLen, hsLen)

	for i := 0; i < hiLen; i++ {
		var h header
		h.parse(r)
		ud.headerIdxs = append(ud.headerIdxs, h)
	}
	for i := 0; i < haLen; i++ {
		var h header
		h.parse(r)
		ud.headerAttrs = append(ud.headerAttrs, h)
	}
	for i := 0; i < hsLen; i++ {
		var h header
		h.parse(r)
		ud.dataStore = append(ud.dataStore, h)
	}
	// fmt.Printf("headerIdxs %+v\n", ud.headerIdxs)
	// fmt.Printf("headerAttrs %+v\n", ud.headerAttrs)
	// fmt.Printf("dataStore %+v\n", ud.dataStore)

	p2Idx := len(data) - 4*5
	p3Idx := len(data) - 4*4
	p2 := BytesToInt(data[p2Idx : p2Idx+4])
	p3 := BytesToInt(data[p3Idx : p3Idx+4])
	// fmt.Println(p2, p3)

	ud.init()
	preOffset := r.Size() - int64(r.Len())
	d := ud.getData(r)
	for i := 0; i < len(d)/2; i++ {

		a, b := d[2*i], d[2*i+1]
		offset := ud.dataStore[ud.keys[0].keyDataIdx].offset + a
		// fmt.Printf("a: %v, b: %v, offset: %v\t", a, b, offset)

		var wordInfo attrWordData
		r.Seek(int64(b), 0)
		wordInfo.parse(r)
		// fmt.Printf("wordInfo: %v", wordInfo)
		// GetWordData
		attrId := ud.keys[0].attrIdx
		dataId := ud.attrs[attrId].dataId
		offset = int(preOffset) + ud.dataStore[dataId].offset + wordInfo.offset
		// fmt.Printf("offset: %v\n", offset)
		// DecryptWordsEx
		word := decryptWordsEx(r, offset, wordInfo.p1, p2, p3)
		ret = append(ret, PyEntry{word, []string{}, wordInfo.freq})
		// fmt.Printf("word: %v\tfreq: %v\n", word, wordInfo.freq)
	}
	return ret
}

func decryptWordsEx(r *bytes.Reader, offset, p1, p2, p3 int) string {
	k1 := (p1 + p2) << 2
	k2 := (p1 + p3) << 2
	xk := (k1 + k2) & 0xFFFF
	r.Seek(int64(offset), 0)
	n := ReadInt(r, 2) / 2
	decWords := make([]byte, 0, 1)
	for i := 0; i < n; i++ {
		shift := p2 % 8
		ch := ReadInt(r, 2)
		dch := (ch<<(16-(shift%8)) | (ch >> shift)) & 0xFFFF
		dch ^= xk
		if dch > 0x10000 {
			print(dch)
		}
		decWords = append(decWords, byte(dch%0x100), byte(dch>>8))
	}
	ret := string(DecUtf16le(decWords))
	return ret
}

type attrWordData struct {
	offset int
	freq   int
	aflag  int
	i8     int
	p1     int
	iE     int
}

func (a *attrWordData) parse(r *bytes.Reader) {
	a.offset = ReadInt(r, 4)
	a.freq = ReadInt(r, 2)
	a.aflag = ReadInt(r, 2)
	a.i8 = ReadInt(r, 4)
	a.p1 = ReadInt(r, 2)
	a.iE = ReadInt(r, 4) // always zero
	_ = ReadInt(r, 4)    // next offset
}

type usrDict struct {
	keys        []key
	attrs       []attr
	aints       []int
	headerIdxs  []header
	headerAttrs []header
	dataStore   []header

	dataTypeSize []int
	attrSize     []int
	baseHashSize []int
	keyHashSize  []int
	aflag        bool
}

var keyDataTypeSize = []int{4, 1, 1, 2, 1, 2, 2, 4, 4, 8, 4, 4, 4, 0, 0, 0}
var dataTypeHashSize = []int{0, 27, 414, 512, -1, -1, 512, 0}

func newUsrDict() *usrDict {
	ud := new(usrDict)
	ud.headerIdxs = make([]header, 0, 1)
	ud.headerAttrs = make([]header, 0, 1)
	ud.dataStore = make([]header, 0, 1)

	ud.dataTypeSize = make([]int, 0, 1)
	// ud.attrSize = make([]int, 0, 1)
	ud.baseHashSize = make([]int, 0, 1)
	ud.keyHashSize = make([]int, 10)
	ud.keyHashSize[0] = 500
	// fmt.Printf("newUsrDict%+v\n", ud)
	return ud
}

func (ud *usrDict) init() {
	ud.attrSize = make([]int, len(ud.attrs))

	for i, k := range ud.keys {
		size := (k.dictType >> 2) & 4
		maskedType := k.dictType & 0xFFFFFF8F
		// hash item
		if ud.keyHashSize[i] > 0 {
			ud.baseHashSize = append(ud.baseHashSize, ud.keyHashSize[i])
		} else {
			ud.baseHashSize = append(ud.baseHashSize, dataTypeHashSize[maskedType])
		}
		// dataType size
		attrCount := ud.attrs[k.attrIdx].count
		// non-attr data size
		nonAttrCount := len(k.dataType) - attrCount
		for j := 0; j < nonAttrCount; j++ {
			if j > 0 || maskedType != 4 {
				size += keyDataTypeSize[k.dataType[i]]
			}
		}
		if k.dictType&0x60 > 0 {
			size += 4
		}
		size += 4
		ud.dataTypeSize = append(ud.dataTypeSize, size)
		// attr data size
		attrSize := 0
		for j := nonAttrCount; j < len(k.dataType); j++ {
			attrSize += keyDataTypeSize[k.dataType[j]]
		}
		if (k.dictType & 0x40) == 0 {
			attrSize += 4
		}
		ud.attrSize[k.attrIdx] = attrSize
		// ???
		if ud.attrs[k.attrIdx].b2 == 0 {
			ud.aflag = true
		}
	}
	// fmt.Printf("init UsrDict%+v\n", ud)
}

func (ud *usrDict) getData(r *bytes.Reader) []int {

	ret := make([]int, 0, 0xff)
	keyId := 0
	theKey := ud.keys[keyId]

	// hashStoreBase := ud.getHashStore(keyId, theKey.dictType&0xFFFFFF8F)
	headerAttr := ud.headerAttrs[theKey.attrIdx]
	var attrCount int

	if headerAttr.usedDataSize == 0 {
		attrCount = headerAttr.dataSize
	} else {
		attrCount = headerAttr.usedDataSize
	}
	hashStoreCount := ud.baseHashSize[keyId]
	// fmt.Println("getData", hashStoreBase, hashStoreCount)

	preOffset := r.Size() - int64(r.Len())
	for i := 0; i < hashStoreCount; i++ {
		r.Seek(preOffset+int64(8*i), 0)
		hashStoreOffset := ReadInt(r, 4)
		hashStoreCount := ReadInt(r, 4)

		// fmt.Printf("hashstore [ offset: {%v}, count: {%v} ]\n", hashStoreOffset, hashStoreCount)
		for j := 0; j < hashStoreCount; j++ {

			attrOffset := int(preOffset) + ud.headerIdxs[keyId].offset + hashStoreOffset + ud.dataTypeSize[keyId]*j
			offset := attrOffset + ud.dataTypeSize[keyId] - 4
			r.Seek(int64(offset), 0)
			offset = ReadInt(r, 4)
			// fmt.Printf("\tattrOffset, %d %d\n", attrOffset, offset)
			for k := 0; k < attrCount; k++ {

				attr2Offset := int(preOffset) + ud.headerAttrs[ud.keys[keyId].attrIdx].offset + offset
				ret = append(ret, attrOffset, attr2Offset)

				offset = attr2Offset + ud.attrSize[theKey.attrIdx] - 4
				r.Seek(int64(offset), 0)
				offset = ReadInt(r, 4)
				// fmt.Printf("\tattr2Offset, %d ,newOffset, %d \n", attr2Offset, offset)
				if offset == 0xFFFFFFFF {
					break
				}
			}
		}
	}
	// fmt.Println(ret)
	return ret
}

func (ud *usrDict) getHashStore(idx, dataType int) int {
	if idx < 0 || dataType > 6 || idx > len(ud.headerIdxs) {
		panic("getHashStore error")
	}
	offset := ud.headerIdxs[idx].offset
	// assert index_offset >= 0
	size := ud.baseHashSize[idx]
	offset = offset - 8*size
	// assert offset >= 0
	return offset
}
