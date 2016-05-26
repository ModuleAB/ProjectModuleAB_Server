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

// ConvertVpcOssAddrToInternal: Aliyun may not support
// OAS pull from OSS with VPC Address
func ConvertVpcOssAddrToInternal(vpcAddr string) string {
	const vpcReg = "vpc100-oss-cn-([a-z]+).aliyuncs.com"
	reg := regexp.MustCompile(vpcReg)
	if reg.MatchString(vpcAddr) {
		region := reg.ReplaceAllString(vpcAddr, "$1")
		return fmt.Sprintf("oss-cn-%s-internal.aliyuncs.com", region)
	}
	return vpcAddr
}
