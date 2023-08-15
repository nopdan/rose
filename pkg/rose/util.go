package rose

import (
	"bytes"
	"fmt"

	"github.com/nopdan/ku"
)

const (
	_u16 uint16 = 0
	_u32 uint32 = 0
)

var (
	ReadUint16 = ku.ReadUint16
	ReadUint32 = ku.ReadUint32

	Encode = ku.Encode
	Decode = ku.Decode
)

func EncodeMust(str, enc string) []byte {
	v, _ := Encode(str, enc)
	return v
}

func DecodeMust(src []byte, enc string) string {
	v, _ := Decode(src, enc)
	return v
}

func PrintInfo(r *bytes.Reader, size uint32, info string) {
	tmp := make([]byte, size)
	r.Read(tmp)
	fmt.Printf("%s%s\n", info, DecodeMust(tmp, "UTF-16LE"))
}
