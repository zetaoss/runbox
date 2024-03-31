package util

import (
	"crypto/md5"
	"fmt"
	"io"
)

func MD5(text string) string {
	hash := md5.New()
	io.WriteString(hash, text)
	md5Bytes := hash.Sum(nil)
	return fmt.Sprintf("%x", md5Bytes)
}
