package encoder

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/cxcn/dtool/pkg/table"
)

func TestShorten(t *testing.T) {
	f, _ := os.Open("../../assets/单字全码表/wubi86.txt")
	defer f.Close()

	wct := make(table.Table, 0, 0xff)
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		wct = append(wct, table.Entry{
			Word: entry[0],
			Code: entry[1],
		})
	}
	Shorten(&wct, "1:3,2:2,4:n")

	var buf bytes.Buffer
	for _, v := range wct {
		buf.WriteString(v.Word)
		buf.WriteByte('\t')
		buf.WriteString(v.Code)
		buf.WriteByte('\n')
	}
	os.WriteFile("own/test_shorten.txt", buf.Bytes(), 0666)
}
