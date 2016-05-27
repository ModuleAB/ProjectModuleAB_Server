package oas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	OasDefaultContentType    = "application/octet-stream"
	OasDefaultSendBufferSize = 8192
	OasDefaultGetBufferSize  = 10 * 1024 * 1024
	OasDefaultProvider       = "OAS"

	OasHttpPort  = 80
	OasHttpsPort = 443

	OasUseHttps = true
	OasNoHttps  = false

	OasUserAgent     = "gooas-OAS Go SDK"
	oasAuthHeaderKey = "Authorization"

	operationFailedFormat = "Operation Failed, %s"
	badHttpResponse       = "Bad HTTP Response"

	OasDescNoContent = ""

	OasJobStatusCodeSucceeded  = "Succeeded"
	OasJobStatusCodeInProgress = "InProgress"
	OasJobStatusCodeFailed     = "Failed"
)

type VaultsList struct {
	Marker    string       `json:"Marker"`
	VaultList []*VaultInfo `json:"VaultList"`
}

type VaultInfo struct {
	CreationDate      string `json:"CreationDate"`
	LastInventoryDate string `json:"LastInventoryDate"`
	NumberOfArchives  int    `json:"NumberOfArchives"`
	SizeInBytes       int    `json:"SizeInBytes"`
	VaultID           string `json:"VaultID"`
	VaultName         string `json:"VaultName"`
}

type JobResult struct {
	Action          string `json:"Action"`
	ArchiveId       string `json:"ArchiveId"`
	ArchiveSize     int    `json:"ArchiveSizeInBytes"`
	TreeEtag        string `json:"TreeEtag"`
	ArchiveTreeEtag string `json:"ArchiveTreeEtag"`
	Completed       bool   `json:"Completed"`
	CompletionDate  string `json:"CompletionDate"`
	CreationDate    string `json:"CreationsDate"`
	InventorySize   int    `json:"InventorySizeInBytes"`
	JobDescription  string `json:"JobDescription"`
	JobId           string `json:"JobId"`
	RetrievalRange  string `json:"RetrievalByteRange"`
	StatusCode      string `json:"StatusCode"`
	StatusMessage   string `json:"StatusMessage"`
}

// 目前只实现需要的功能，其他功能需要时才实现
// （不知道Go对大文件上传的支持怎样）
type OasClient struct {
	host      string
	apiKey    string
	apiSecret string

	h *http.Client
}

func NewOasClient(host, apikey, secret string, port int, security bool) *OasClient {
	o := new(OasClient)
	o.apiKey = apikey
	o.apiSecret = secret
	if security || port == 443 {
		o.host = fmt.Sprintf("https://%s:%d", host, port)
	} else {
		o.host = fmt.Sprintf("http://%s:%d", host, port)
	}
	o.h = new(http.Client)
	return o
}

func (o *OasClient) getResource(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}

	tmpHeaders := make(map[string]interface{})
	for k, v := range params {
		tmpK := strings.TrimSpace(strings.ToLower(k))
		tmpHeaders[tmpK] = v
	}

	overrideResponseList := []string{
		"limit", "marker", "response-content-type", "response-content-language",
		"response-cache-control", "logging", "response-content-encoding",
		"acl", "uploadId", "uploads", "partNumber", "group",
		"delete", "website", "location", "objectInfo",
		"response-expires", "response-content-disposition"}
	sort.Strings(overrideResponseList)

	resource := ""
	separator := "?"
	for _, i := range overrideResponseList {
		if _, ok := tmpHeaders[strings.ToLower(i)]; ok {
			resource = fmt.Sprintf("%s%s%s", resource, separator, i)
			tmpKey := tmpHeaders[strings.ToLower(i)]
			if tmpKey != "" {
				resource = fmt.Sprintf("%s=%v", resource, tmpKey)
			}
			separator = "&"
		}
	}
	return resource
}

// 底层方法，用于向OAS发送请求
// 其中headers, body, params可以为nil
func (o *OasClient) httpRequest(method, url string, headers http.Header,
	body *bytes.Buffer, params map[string]interface{}) (*http.Response, error) {
	if headers == nil {
		headers = make(http.Header)
	}
	headers.Set("User-Agent", OasUserAgent)
	headers.Set("Host", o.host)
	headers.Set("Date", strings.Replace(time.Now().UTC().Format(time.RFC1123),
		"UTC", "GMT", 1))
	headers.Set("x-oas-version", "0.2.5")
	if body.Len() > 0 {
		headers.Set("Content-Length", fmt.Sprint(body.Len()))
	}

	resource := fmt.Sprintf("%s%s", url, o.getResource(params))
	if params != nil {
		url = appendParam(url, params)
	}
	headers.Set(oasAuthHeaderKey,
		o.createSignForNormalAuth(method, headers, resource))

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", o.host, url), body)
	if err != nil {
		return nil, err
	}

	req.Header = headers
	return o.h.Do(req)
}

func (o *OasClient) createSignForNormalAuth(method string, headers http.Header,
	resource string) string {
	var res []string
	authValue := fmt.Sprintf("%s %s:%s", OasDefaultProvider, o.apiKey,
		getAssign(o.apiSecret, method, headers, resource, &res))
	return authValue
}

// 获取Vault列表
// 当limit值为-1时，limit参数不启用
// 当marker值为""时，marker参数不启用
func (o *OasClient) ListVaults(limit int, marker string) (requestId string,
	v *VaultsList, err error) {

	defer func() {
		x := recover()
		if x != nil {
			err = fmt.Errorf("%v", x)
		}
	}()

	url := "/vaults"
	method := "GET"

	params := make(map[string]interface{})
	if limit > -1 {
		params["limit"] = limit
	}
	if marker != "" {
		params["marker"] = marker
	}
	if len(params) == 0 {
		params = nil
	}

	r, err := o.httpRequest(method, url, nil, new(bytes.Buffer), params)
	if err != nil {
		return
	}
	requestId = r.Header.Get("x-oas-request-id")
	if err = checkResponse(r, http.StatusOK); err != nil {
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	v = new(VaultsList)
	v.VaultList = make([]*VaultInfo, 0)

	err = json.Unmarshal(b, v)
	if err != nil {
		return
	}
	err = nil
	return
}

func (o *OasClient) ArchiveToOas(vaultID, ossHost, bucket,
	path, desc string) (requestId, jobId string, err error) {
	defer func() {
		x := recover()
		if x != nil {
			err = fmt.Errorf("%v", x)
		}
	}()
	url := fmt.Sprintf("/vaults/%s/jobs", vaultID)
	method := "POST"
	reqJson := &struct {
		Type    string `json:"Type"`
		Desc    string `json:"Description"`
		OSSHost string `json:"OSSHost"`
		Bucket  string `json:"Bucket"`
		Object  string `json:"Object"`
	}{
		Type:    "pull-from-oss",
		Desc:    desc,
		OSSHost: ossHost,
		Bucket:  bucket,
		Object:  path,
	}
	b, err := json.Marshal(reqJson)
	if err != nil {
		return
	}
	body := bytes.NewBuffer(b)
	r, err := o.httpRequest(method, url, nil, body, nil)
	if err != nil {
		return
	}
	requestId = r.Header.Get("x-oas-request-id")
	if err = checkResponse(r, http.StatusAccepted); err != nil {
		return
	}

	jobId = r.Header.Get("x-oas-job-id")
	err = nil
	return
}

func (o *OasClient) RecoverToOss(vaultID, archiveId, ossHost,
	bucket, path, desc string) (requestId, jobId string, err error) {
	defer func() {
		x := recover()
		if x != nil {
			err = fmt.Errorf("%v", x)
		}
	}()
	url := fmt.Sprintf("/vaults/%s/jobs", vaultID)
	method := "POST"
	reqJson := &struct {
		Type      string `json:"Type"`
		Desc      string `json:"Description"`
		OSSHost   string `json:"OSSHost"`
		Bucket    string `json:"Bucket"`
		ArchiveId string `json:"ArchiveId"`
		Object    string `json:"Object"`
	}{
		Type:      "push-to-oss",
		Desc:      desc,
		OSSHost:   ossHost,
		Bucket:    bucket,
		Object:    path,
		ArchiveId: archiveId,
	}
	b, err := json.Marshal(reqJson)
	if err != nil {
		return
	}
	body := bytes.NewBuffer(b)
	r, err := o.httpRequest(method, url, nil, body, nil)
	if err != nil {
		return
	}
	requestId = r.Header.Get("x-oas-request-id")
	if err = checkResponse(r, http.StatusAccepted); err != nil {
		return
	}

	jobId = r.Header.Get("x-oas-job-id")
	err = nil
	return
}

func (o *OasClient) GetJobInfo(vaultID, jobId string) (requestId string,
	jr *JobResult, err error) {
	defer func() {
		x := recover()
		if x != nil {
			err = fmt.Errorf("%v", x)
		}
	}()
	url := fmt.Sprintf("/vaults/%s/jobs/%s", vaultID, jobId)
	method := "GET"
	r, err := o.httpRequest(method, url, nil, new(bytes.Buffer), nil)
	if err != nil {
		return
	}
	requestId = r.Header.Get("x-oas-request-id")

	if err = checkResponse(r, http.StatusOK); err != nil {
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	jr = new(JobResult)
	err = json.Unmarshal(b, jr)
	if err != nil {
		return
	}
	err = nil
	return
}

func (o *OasClient) DeleteArchive(vaultID, archiveId string) (requestId,
	jobId string, err error) {
	defer func() {
		x := recover()
		if x != nil {
			err = fmt.Errorf("%v", x)
		}
	}()
	url := fmt.Sprintf("/vaults/%s/archives/%s", vaultID, archiveId)
	method := "DELETE"
	r, err := o.httpRequest(method, url, nil, new(bytes.Buffer), nil)
	if err != nil {
		return
	}
	requestId = r.Header.Get("x-oas-request-id")
	if err = checkResponse(r, http.StatusNoContent); err != nil {
		return
	}
	jobId = r.Header.Get("x-oas-job-id")
	err = nil
	return
}
