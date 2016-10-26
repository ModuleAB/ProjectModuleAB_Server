/*ModuleAB common/oss.go -- Aliyun OSS instance.
 * Copyright (C) 2016 TonyChyi <tonychee1989@gmail.com>
 * License: GPL v3 or later.
 */

package common

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego"
)

type OssClient struct {
	*oss.Client
}

func NewOssClient(endpoint string) (*OssClient, error) {
	if !strings.HasPrefix(
		"http://",
		strings.ToLower(endpoint),
	) {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}

	var err error
	o := new(OssClient)
	o.Client, err = oss.New(
		endpoint,
		beego.AppConfig.String("aliapi::apikey"),
		beego.AppConfig.String("aliapi::secret"),
	)
	return o, err
}

// ConvertOssAddrToInternal: Aliyun may not support
// OAS pull from OSS with VPC Address
func ConvertOssAddrToInternal(addr string) string {
	const Reg = "(?:vpc100-)?oss-cn-([a-z0-9-]+).aliyuncs.com"
	reg := regexp.MustCompile(Reg)

	if reg.MatchString(addr) {
		region := reg.ReplaceAllString(addr, "$1")
		if !regexp.MustCompile("-internal").MatchString(region) {
			return fmt.Sprintf("oss-cn-%s-internal.aliyuncs.com", region)
		}
	}
	return addr
}
