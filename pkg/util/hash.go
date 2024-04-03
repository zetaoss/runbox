package util

import (
	"bytes"
	"crypto/sha256"

	"github.com/google/uuid"
)

func NewHash(length int) string {
	return Hash(uuid.New().String()+uuid.New().String(), length)
}

func Hash(s string, n int) string {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte(s))
	hashBytes := hasher.Sum(nil)

	baseChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var baseHash bytes.Buffer
	for _, b := range hashBytes {
		baseHash.WriteByte(baseChars[b%36])
	}
	var result bytes.Buffer
	cnt := n / baseHash.Len()
	remainder := n % baseHash.Len()
	for i := 0; i < cnt; i++ {
		result.Write(baseHash.Bytes())
	}
	result.WriteString(baseHash.String()[:remainder])
	return result.String()
}
