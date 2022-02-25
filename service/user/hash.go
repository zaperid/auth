package user

import (
	"crypto/sha1"
	"encoding/hex"
)

func hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	rawHash := hasher.Sum(nil)
	hashHex := hex.EncodeToString(rawHash)

	return hashHex
}
