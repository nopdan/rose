package wubi

import (
	"bytes"
	"fmt"
	"io"
)

// TODO:
type Kafan struct{ Template }

func NewKafan() *Kafan {
	return &Kafan{}
}
func (k *Kafan) Unmarshal(r *bytes.Reader) []*Entry {
	wl := make([]*Entry, 0, 0xff)
	r.Seek(0x48, io.SeekStart)
	check := make([]byte, 0x10)
	r.Read(check)
	check1 := string(check)
	r.Read(check)
	check2 := string(check)
	if check1 != "ProtoDict1" || check2 != "" {
		fmt.Println("格式错误！")
	}
	return wl
}

func (k *Kafan) Marshal() []byte {
	return nil
}
