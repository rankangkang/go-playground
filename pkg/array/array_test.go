package array_test

import (
	"fmt"
	"playground/pkg/array"
	"testing"
)

func TestArray(t *testing.T) {
	arr := []int{1, 2, 3, 4}
	array.ForEach(arr, func(item int, i int) {
		fmt.Println(i)
	})

	nextArr := array.Map(arr, func(item int, i int) int {
		return item + 1
	})
	fmt.Println(nextArr)

	m := array.Array2Map(nextArr, func(item int, i int) string {
		return fmt.Sprintf("%v", item)
	})
	fmt.Println(m)

	v, ok := m["2"]
	if ok {
		fmt.Println(v)
	}

	fmt.Println(array.Contains(arr, 1))
}
