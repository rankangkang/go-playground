package util

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// 转换 struct，将 src 与 dst struct 中相同的属性，从 src 复制到 dst，二者都要传指针
func StructTransfer(src interface{}, dst interface{}) error {
	if src == nil {
		return fmt.Errorf("src struct is nil")
	}

	if dst == nil {
		return fmt.Errorf("dst struct is nil")
	}

	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return fmt.Errorf("dst struct is not a pointer")
	}

	bs, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("src marshal err=%v", err)
	}
	err = json.Unmarshal(bs, dst)
	if err != nil {
		return fmt.Errorf("src unmarshal err=%v", err)
	}

	return nil
}
