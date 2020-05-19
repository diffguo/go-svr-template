package controller

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
)

const HmacKey = "973dfe463ec85785f5f95af5ba3906eedb2d931c24e69824a89ea65dba98763b"

func Hmac4Password(password string) string {
	h := hmac.New(sha256.New, []byte(HmacKey))
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}
