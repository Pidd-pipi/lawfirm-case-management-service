package util

// Unique 返回去重后的 uint64 切片。
func Unique(s []uint64) []uint64 {
	out := s[:0]
	for _, v := range s {
		if !containsU64(out, v) {
			out = append(out, v)
		}
	}
	return out
}

// Filter 返回满足 keep 条件的元素构成的新切片。
func Filter(s []uint64, keep func(uint64) bool) []uint64 {
	out := s[:0]
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

func containsU64(list []uint64, v uint64) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
