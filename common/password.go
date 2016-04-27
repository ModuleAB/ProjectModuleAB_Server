package common

import (
	"crypto/sha1"
	"encoding/base64"
	"strings"
)

func EncryptPassword(password string) string {
	hash := sha1.New()
	b := strings.NewReader(password)
	b.WriteTo(hash)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
