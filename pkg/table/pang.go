package table

import (
	"bytes"
	"os"
)

type Pang struct{}

func (Pang) Parse(filename string) Table {
	data, _ := os.ReadFile(filename)
	r := bytes.NewReader(data)
	ret := make(Table, 0, r.Len()>>8)
	var tmp []byte

	r.Seek(0x16, 0) // 从 0x16 开始读
	for r.Len() > 10 {
		var flag = true
		for flag {
			b, _ := r.ReadByte()
			switch b {
			case 0x00, 0xff:
			case 0x02:
				tmp = make([]byte, 4)
				r.Read(tmp)
				if tmp[3] != 0x00 {
					r.Seek(-5, 1)
					flag = false
				}
			case 0x01, 0x06:
				r.Seek(-1, 1)
				flag = false
			default:
				flag = false
			}
		}

		first, _ := r.ReadByte()
		var word []rune

		switch first {
		case 0x06:
		default:
			r.Seek(4, 1) // 前 4 个字节不用。
		}

		r.Seek(4, 1)
		wordCount := ReadUint32(r)
		word = make([]rune, 0)
		for i := 0; i < int(wordCount); i++ {
			b, _ := r.ReadByte()
			if b > 0xF0 {
				r.Seek(-1, 1)
				break
			} else {
				r.Seek(-1, 1)
			}

			tmp1 := make([]byte, 2)
			r.Read(tmp1)
			if bytes.Equal(tmp1, []byte{0xc2, 0xb7}) {
				continue
			} else {
				r.Seek(-2, 1)
			}

			tmp2, _, err := r.ReadRune()
			if err != nil {
				r.Seek(-1, 1)
				break
			}
			if len(string(tmp2)) == 4 {
				i++
			}
			word = append(word, tmp2)
		}

		r.Seek(4, 1)
		codeLen := ReadUint32(r)
		tmp = make([]byte, codeLen)
		r.Read(tmp)

		pos := ReadUint32(r)
		ret = append(ret, Entry{string(word), string(tmp), int(pos)})
		// fmt.Printf("词:%s, 编码:%s, 位置:%d\n", string(word), string(tmp), pos)
	}
	return ret
}
