package oas

import (
	"encoding/json"
	"testing"
)

const (
	host   = "cn-hangzhou.oas.aliyuncs.com"
	key    = ""
	secret = ""
)

func Test(t *testing.T) {
	client := NewOasClient(host, key, secret, 80, false)
	t.Log("Host:", client.host)
	id, lists, err := client.ListVaults(-1, "")
	if err != nil {
		t.Error("Get Error at ListVaults():", err)
	}
	t.Log("Request ID:", id)
	b, err := json.Marshal(lists)
	if err != nil {
		t.Error("Get Error at json:", err)
	}
	t.Log(string(b))
}
