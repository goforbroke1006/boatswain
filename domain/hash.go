package domain

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSHA256(phrase []byte) string {
	h := sha256.New()
	h.Write(phrase)
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
