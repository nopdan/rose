package baidu_bak

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/util"
)

const MASK = 0x2D382324

var utf16 = util.NewEncoding("UTF-16LE")

var (
	TABLE        = []byte("qogjOuCRNkfil5p4SQ3LAmxGKZTdesvB6z_YPahMI9t80rJyHW1DEwFbc7nUVX2-")
	DECODE_TABLE [256]byte
)

func init() {
	for i, b := range TABLE {
		DECODE_TABLE[b] = byte(i)
	}
}

// BaiduBak 百度拼音备份.bin格式
type BaiduBak struct {
	model.BaseFormat
}

// New 创建百度词库格式处理器
func New() *BaiduBak {
	return &BaiduBak{
		BaseFormat: model.BaseFormat{
			ID:          "baidu_bak",
			Name:        "百度拼音备份.bin",
			Type:        model.FormatTypePinyin,
			Extension:   ".bin",
			Description: "百度拼音备份二进制格式",
		},
	}
}

const (
	T_cnword = iota
	T_enword
	T_sysusrword
)

func (f *BaiduBak) Import(src model.Source) ([]*model.Entry, error) {
	r, err := model.NewReaderFromSource(src)
	if err != nil {
		return nil, err
	}

	entries := make([]*model.Entry, 0, r.Size()>>8)

	// 跳过开头的 FF FE
	r.Seek(2, io.SeekStart)

	var currentType int

	for r.Len() > 8 {
		lineBytes := make([]byte, 0, 20)
		for {
			// 每次读取两个字节，添加到 line，直到遇到换行符 0x0A 0x00
			b := make([]byte, 2)
			n, err := r.Read(b)
			if err != nil {
				os.Exit(1)
			}
			if b[0] == 0x0A && b[1] == 0x00 {
				break
			}
			lineBytes = append(lineBytes, b[:n]...)
		}
		lineBytes, err = decode(lineBytes)
		if err != nil {
			continue
		}
		line := utf16.Decode(lineBytes)
		f.Infof("%s\n", line)
		switch line {
		case "<cnword>":
			currentType = T_cnword
			continue
		case "<enword>":
			currentType = T_enword
			continue
		case "<sysusrword>":
			currentType = T_sysusrword
			continue
		}
		switch currentType {
		case T_cnword:
			// example: 词库(ci|ku) 2 24 1703756732 N N
			items := strings.Split(line, " ")
			if len(items) < 3 {
				continue
			}
			rank, _ := strconv.Atoi(items[1])
			w := items[0]                       // 词库(ci|ku)
			w = strings.TrimSuffix(w, ")")      // 词库(ci|ku
			items = strings.Split(w, "(")       // [词库, ci|ku]
			word := items[0]                    // 词库
			pys := strings.Split(items[1], "|") // [ci, ku]
			entry := model.NewEntry(word).WithMultiCode(pys...).WithRank(rank)
			entries = append(entries, entry)
		case T_enword:
			items := strings.Split(line, " ")
			if len(items) < 2 {
				continue
			}
			entry := model.NewEntry(items[0])
			entries = append(entries, entry)
		case T_sysusrword:
		}

	}
	return entries, nil
}

func decode(data []byte) ([]byte, error) {
	if len(data)%4 != 2 {
		return nil, errors.New("invalid data length")
	}

	base64Remainder := data[len(data)-2] - 65
	if base64Remainder > 2 || data[len(data)-1] != 0 {
		return nil, errors.New("invalid padding")
	}

	// 映射魔改过的 base64 编码表
	for i := 0; i < len(data)-2; i++ {
		data[i] = DECODE_TABLE[data[i]]
	}

	result := make([]byte, 0, len(data)/4*3)
	for i := 0; i < len(data)-2; i += 4 {
		highBits := data[i+3]
		result = append(result, byte(data[i]|(highBits&0b110000)<<2))
		result = append(result, byte(data[i+1]|(highBits&0b1100)<<4))
		result = append(result, byte(data[i+2]|(highBits&0b11)<<6))
	}

	if base64Remainder != 0 {
		for i := byte(0); i < 3-base64Remainder; i++ {
			if result[len(result)-1] != 0 {
				return nil, errors.New("invalid padding")
			}
			result = result[:len(result)-1]
		}
	}

	for i := 0; i < len(result); i += 4 {
		if i+4 > len(result) {
			data := make([]byte, 4)
			copy(data, result[i:])
			chunk := MASK ^ binary.LittleEndian.Uint32(data)
			binary.LittleEndian.PutUint32(data, chunk)
			copy(result[i:], data[:len(result)-i])
			break
		}
		chunk := MASK ^ binary.LittleEndian.Uint32(result[i:i+4])
		chunk = (chunk&0x1FFFFFFF)<<3 | chunk>>29
		binary.LittleEndian.PutUint32(result[i:i+4], chunk)
	}
	return result, nil
}
