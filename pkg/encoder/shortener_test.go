package encoder

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/cxcn/dtool/pkg/util"
)

func TestShorten(t *testing.T) {
	f, _ := os.Open("own/星辰星笔全码.txt")
	defer f.Close()

	wct := make(WcTable, 0, 0xff)
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) < 2 {
			continue
		}
		wct = append(wct, WordCode{
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
	ioutil.WriteFile("own/星辰星笔全码_shorten.txt", buf.Bytes(), 0777)
}
