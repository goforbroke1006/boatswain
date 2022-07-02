package internal

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSHA256(phrase string) string {
	h := sha256.New()
	h.Write([]byte(phrase))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
