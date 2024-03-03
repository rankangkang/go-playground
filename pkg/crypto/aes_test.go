package crypto_test

import (
	"playground/pkg/crypto"
	"testing"
)

func TestAes(t *testing.T) {
	key := "0123456789abcdef"
	iv := "0123456789012345"
	c := crypto.NewCbc(key, iv)
	err := c.EncryptFileFull("../../test/test.txt", "../../test/test-out.txt")
	if err != nil {
		t.Error(err)
	}

	err = c.DecryptFileFull("../../test/test-out.txt", "../../test/test-out2.txt")
	if err != nil {
		t.Error(err)
	}

	// var m = map[string]any{}
	// err = fs.ReadJson("../../test/test-out2.txt", &m)
	// if err != nil {
	// 	t.Error(err)
	// }

	// fmt.Println(m)
}
