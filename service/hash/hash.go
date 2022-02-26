package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	hashByte := hasher.Sum(nil)
	hashStr := hex.EncodeToString(hashByte)

	return hashStr
}
