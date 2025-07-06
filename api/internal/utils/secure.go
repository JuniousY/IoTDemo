package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func SecureRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

func HmacSha256Pass(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
