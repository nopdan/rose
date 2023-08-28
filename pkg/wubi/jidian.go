package wubi

import (
	"bufio"
	"bytes"
	"sort"
	"strings"
)

type Jidian struct{ Template }

func init() {
	FormatList = append(FormatList, NewJidian())
}
func NewJidian() *Jidian {
	f := new(Jidian)
	f.Name = "极点码表"
	f.ID = "jidian"
	f.CanMarshal = true
	f.HasRank = true
	return f
}

func (Jidian) Unmarshal(r *bytes.Reader) []*Entry {
	di := make([]*Entry, 0, r.Size()>>8)

	scan := bufio.NewScanner(r)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), " ")
		if len(entry) < 2 {
			continue
		}
		for i := 1; i < len(entry); i++ {
			di = append(di, &Entry{
				Word: entry[i],
				Code: entry[0],
				Rank: i,
			})
		}
	}
	return di
}

func (Jidian) Marshal(di []*Entry, hasRank bool) []byte {
	var buf bytes.Buffer
	buf.Grow(len(di))

	sort.SliceStable(di, func(i, j int) bool {
		return di[i].Code < di[j].Code
	})
	// 记录上一个编码
	lastCode := ""
	for i, v := range di {
		if lastCode != v.Code {
			if i != 0 {
				buf.WriteByte('\r')
				buf.WriteByte('\n')
			}
			buf.WriteString(v.Code)
			buf.WriteByte('\t')
			buf.WriteString(v.Word)
			lastCode = v.Code
			continue
		}
		buf.WriteByte(' ')
		buf.WriteString(v.Word)
		lastCode = v.Code
	}
	return buf.Bytes()
}
