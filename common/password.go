/*ModuleAB common/password.go -- encrypt password with sha1.
 * Copyright (C) 2016 TonyChyi <tonychee1989@gmail.com>
 * License: GPL v3 or later.
 */

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
