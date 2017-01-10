/*ModuleAB common/auth.go -- define auth algorithm
 * Copyright (C) 2016 TonyChyi <tonychee1989@gmail.com>
 * License: GPL v3 or later.
 */

package common

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

//AuthExpireDuration is 10 minutes
const AuthExpireDuration = 10 * time.Minute

//Auth whether http request is valid.
func AuthWithKey(ctx *context.Context) error {
	key := beego.AppConfig.String("loginkey")
	sTime := ctx.Input.Header("Date")
	pTime, err := time.Parse(time.RFC1123, sTime)
	if err != nil {
		return err
	}
	if time.Now().UTC().Sub(pTime) > AuthExpireDuration ||
		time.Now().UTC().Sub(pTime) < -AuthExpireDuration {
		return fmt.Errorf("Client time is out of server time")
	}
	sign := ctx.Input.Header("Signature")
	h := hmac.New(sha1.New, []byte(key))
	beego.Debug("Got URL:", ctx.Input.URL())
	beego.Debug("Got date:", sTime)
	h.Write(
		[]byte(
			fmt.Sprintf(
				"%s%s",
				sTime,
				ctx.Input.URL(),
			),
		),
	)
	b := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if sign != b {
		return fmt.Errorf("Bad signature.")
	}
	return nil
}
