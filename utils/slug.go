package utils

import (
	"crypto/rand"
	"math/big"
)

const slugChars = "abcdefghijklmnopqrstuvwxyz0123456789"

func PublicSlug() string {
	b := make([]byte, 8)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(slugChars))))
		b[i] = slugChars[n.Int64()]
	}
	return string(b)
}
