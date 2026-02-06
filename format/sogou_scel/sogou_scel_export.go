package sogou_scel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

func (f *SogouScel) Export(entries []*model.Entry, w io.Writer) error {
	if len(syllMap) == 0 {
		syllMap = list2map(pyList)
		pyBytes = list2bytes(pyList)
	}

	var buf bytes.Buffer
	buf.Write(make([]byte, 0x1540))
	buf.Write(pyBytes)
	cCount, wCount, cSize, wSize := 0, 0, 0, 0

	//!TODO 暂未使用 合并相同拼音的词条
	//!TODO 短语暂未完全解析
	// 示例词
	examples := make([]string, 0)

	for _, entry := range entries {
		buf.Write([]byte{0x01, 0x00})

		sylls, word := entry.Code.Strings(), entry.Word
		// 拼音占用字节数
		buf.Write(util.To2Bytes(len(sylls) * 2))
		for _, s := range sylls {
			idx := index(s)
			buf.Write(util.To2Bytes(idx))
		}

		cSize += len(sylls)*2 + 2
		cCount++

		// 添加示例词
		if len(examples) < 6 {
			examples = append(examples, word)
		}
		b := utf16.Encode(word)
		buf.Write(util.To2Bytes(len(b)))
		buf.Write(b)
		buf.Write([]byte{0x0A, 0x00})
		buf.Write(util.To4Bytes(entry.Frequency))
		buf.Write(make([]byte, 6))

		wSize += len(b) + 2
		wCount++
	}

	b := buf.Bytes()
	copy(b, []byte{0x40, 0x15, 0x00, 0x00, 0xD2, 0x6D, 0x53, 0x01, 0x01, 0x00, 0x00, 0x00})

	// 校验和
	chkrd := bytes.NewReader(b[0x1540:])
	chksum := checkSumStream(chkrd)
	for i, v := range chksum {
		tmp := util.To4Bytes(v)
		copy(b[i*4+0xC:], tmp)
	}

	// 生成 ID
	rand_num := rand.Uint32()
	str := fmt.Sprintf("L%d", uint16(rand_num))
	copy(b[0x1C:], utf16.Encode(str))

	// 时间戳
	now := uint32(time.Now().Unix())
	binary.LittleEndian.PutUint32(b[0x11C:], now)

	copy(b[0x120:], util.To4Bytes(cCount))
	copy(b[0x124:], util.To4Bytes(wCount))
	copy(b[0x128:], util.To4Bytes(cSize))
	copy(b[0x12C:], util.To4Bytes(wSize))

	copy(b[NAME_OFFSET:], utf16.Encode("词库名称"))
	copy(b[LOCATION_OFFSET:], utf16.Encode("本地"))
	copy(b[DESC_OFFSET:], utf16.Encode("由 scel-maker 生成的细胞词库"))
	copy(b[EXAMPLE_OFFSET:], utf16.Encode(strings.Join(examples, "   ")))

	_, err := w.Write(b)
	return err
}
