package util

func RmRepeat[T comparable](sli []T) []T {
	ret := make([]T, 0, len(sli))
	for _, v := range sli {
		if !Contain(ret, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func Contain[T comparable](sli []T, el T) bool {
	for _, v := range sli {
		if v == el {
			return true
		}
	}
	return false
}
