package array

// array 转 map
func Array2Map[T any, K string | int | int64](slice []T, getKey func(T, int) K) map[K]T {
	result := map[K]T{}
	for i, item := range slice {
		key := getKey(item, i)
		result[key] = item
	}

	return result
}
