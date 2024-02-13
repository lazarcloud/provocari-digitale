package auth

import (
	"crypto/sha1"
	"encoding/hex"
)

func HashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash

}

func CompareHash(s string, hash string) bool {
	return HashString(s) == hash
}
