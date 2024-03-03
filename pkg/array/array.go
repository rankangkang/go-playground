package array

// forEach
func ForEach[T any](slice []T, action func(T, int)) {
	for i, item := range slice {
		action(item, i)
	}
}

// map
func Map[T1 any, T2 any](slice []T1, mapper func(T1, int) T2) []T2 {
	mapped := make([]T2, 0, len(slice))
	for i, item := range slice {
		mapped = append(mapped, mapper(item, i))
	}
	return mapped
}

// filter
func Filter[T any](slice []T, condition func(T, int) bool) []T {
	var filtered []T
	for i, item := range slice {
		if condition(item, i) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// reduce
func Reduce[T any](slice []T, reduce func(cur T, item T, i int) T, init T) T {
	for i := 0; i < len(slice); i++ {
		init = reduce(init, slice[i], i)
	}
	return init
}

// every
func Every[T any](slice []T, condition func(T, int) bool) bool {
	for i, item := range slice {
		if !condition(item, i) {
			return false
		}
	}
	return true
}

// some
func Some[T any](slice []T, condition func(T, int) bool) bool {
	for i, item := range slice {
		if condition(item, i) {
			return true
		}
	}
	return false
}

// find
func Find[T any](slice []T, condition func(T, int) bool) (T, bool) {
	for i, item := range slice {
		if condition(item, i) {
			return item, true
		}
	}
	var t T
	return t, false
}

// findIndex
func FindIndex[T any](slice []T, condition func(T, int) bool) int {
	for i, item := range slice {
		if condition(item, i) {
			return i
		}
	}
	return -1
}

// indexOf
func IndexOf[T comparable](slice []T, item T) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

// lastIndexOf
func LastIndexOf[T comparable](slice []T, item T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

// contains
func Contains[T comparable](slice []T, target T) bool {
	for _, val := range slice {
		var v T = val
		if target == v {
			return true
		}
	}

	return false
}

// 两切片是否相交
func Intersects[T comparable](s, t []T) bool {
	for _, item := range s {
		if Contains(t, item) {
			return true
		}
	}

	return false
}

// 两切片内容是否完全相同，包括顺序
func FullyEqual[T comparable](s, t []T) bool {
	if len(s) != len(t) {
		return false
	}

	length := len(s)
	for i := 0; i < length; i++ {
		if s[i] != t[i] {
			return false
		}
	}

	return true
}

// 切片内容是否相同，不关心顺序
func Equal[T comparable](s, t []T) bool {
	if len(s) != len(t) {
		return false
	}

	cm1 := make(map[T]int)
	cm2 := make(map[T]int)
	for _, v := range s {
		cm1[v]++
	}
	for _, v := range t {
		cm2[v]++
	}

	if len(cm1) != len(cm2) {
		return false
	}

	for k, freq1 := range cm1 {
		freq2, ok := cm2[k]
		if !ok || freq1 != freq2 {
			return false
		}
	}

	return true
}

// 数组去重
func Deduplicate[T comparable](src []T) []T {
	var seen = make(map[T]bool)
	var dst = []T{}

	for _, v := range src {
		if !seen[v] {
			dst = append(dst, v)
			seen[v] = true
		}
	}

	return dst
}
