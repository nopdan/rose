package util

import (
	"fmt"
	"testing"
)

func TestProduct(t *testing.T) {
	sli := [][]byte{{'a', 'b'}, {'c'}, {'d', 'e', 'f'}}
	new := Product(sli)
	for _, v := range new {
		fmt.Println(string(v))
	}

	sli = [][]byte{{'a', 'b'}}
	new = Product(sli)
	for _, v := range new {
		fmt.Println(string(v))
	}
}
