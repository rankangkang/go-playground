package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSha256Sum(content string) string {
	h := sha256.New()
	h.Write([]byte(content))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
