package encoder

// 通过词表+单字表生成词的拼音
func GetPinyin(word string) []string {
	chars := []rune(word)
	ret := make([]string, 0, len(chars))

	for i := 0; i < len(chars); {
		// 匹配到词
		flag := false
		for j := len(chars); i+1 < j; j-- {
			if pys, ok := WordPinyinMap[string(chars[i:j])]; ok {
				ret = append(ret, pys...)
				i = j
				flag = true
			}
		}
		if flag {
			continue
		}
		if py, ok := CharYinjieMap[chars[i]]; ok {
			ret = append(ret, py[0])
		} else {
			ret = append(ret, "")
		}
		i++
	}
	return ret
}

// 通过单字表生成词的拼音
func GetPyByChar(word string) []string {
	chars := []rune(word)
	ret := make([]string, 0, len(chars))
	for _, char := range chars {
		if py, ok := CharYinjieMap[char]; ok {
			ret = append(ret, py[0])
		} else {
			ret = append(ret, "")
		}
	}
	return ret
}
