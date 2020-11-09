package util

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

var hmacSecret []byte

// EncodeMD5 md5 encryption
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

// Hmac  md5 with specified key
func Hmac(data string) string {
	hash := hmac.New(md5.New, hmacSecret)
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}