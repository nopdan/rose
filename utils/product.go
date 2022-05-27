package util

type A [][]byte

// 笛卡尔积
func Product(a A) []string {
	if len(a) <= 1 {
		return []string{}
	}
	res := make(A, 0, len(a))
	for _, v := range a[0] {
		res = append(res, []byte{v})
	}
	for i := 1; i < len(a); i++ {
		res = productOne(res, a[i])
	}
	ret := make([]string, 0, len(res))
	for _, v := range res {
		ret = append(ret, string(v))
	}
	return ret
}

func productOne(a A, b []byte) A {
	ret := make([]string, 0, len(a)*len(b))
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			ret = append(ret, string(append(a[i], b[j])))
		}
	}
	tmp := make(A, 0, len(a)*len(b))
	for _, v := range ret {
		tmp = append(tmp, []byte(v))
	}
	return tmp
}
