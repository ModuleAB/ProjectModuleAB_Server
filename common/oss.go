package common

import (
	"fmt"
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
