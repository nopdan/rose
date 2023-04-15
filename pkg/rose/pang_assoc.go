package rose

import "bytes"

type PangAssoc struct{}

func (PangAssoc) Parse(b []byte) WubiTable {
	r := bytes.NewReader(b)
	ret := make(WubiTable, 0, len(b)>>8)
	var tmp []byte

	readWord := func(word *[]rune, r *bytes.Reader) {
		wordCount := ReadUint32(r)
		for i := 0; i < int(wordCount); i++ {
			b, _ := r.ReadByte()
			if b > 0xF0 {
				r.Seek(-1, 1)
				break
			} else {
				r.Seek(-1, 1)
			}
			tmp2, _, err := r.ReadRune()
			if err != nil {
				r.Seek(-1, 1)
				break
			}
			if len(string(tmp2)) == 4 {
				i++
			}
			*word = append(*word, tmp2)
		}
	}

	r.Seek(0x4, 0) // 从 0x16 开始读
	for r.Len() > 10 {

		var word []rune
		r.Seek(5, 1)
		readWord(&word, r)

		r.Seek(4, 1)
		codeLen := ReadUint32(r)
		tmp = make([]byte, codeLen)
		r.Read(tmp)

		r.Seek(9, 1)
		word = append(word, '#')
		readWord(&word, r)
		r.Seek(1, 1)

		ret = append(ret, &WubiEntry{string(word), string(tmp), 1})
		// fmt.Printf("词:%s, 编码:%s\n", string(word), string(tmp))
	}
	return ret
}
