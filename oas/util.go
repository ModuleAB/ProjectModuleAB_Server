package oas

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const (
	OasDefineHeaderPrefix = "x-oas-"
)

type ErrorMsg struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *ErrorMsg) String() string {
	return fmt.Sprintf(`Code: "%s", Message: "%s", Type: "%s"`,
		e.Code, e.Message, e.Type)
}

func appendParam(Url string, params map[string]interface{}) string {
	l := make([]string, 0)
	for k, v := range params {
		k = strings.Replace(k, "_", "-", -1)
		if k == "maxkeys" {
			k = "max-keys"
		}

		nv := fmt.Sprint(v)
		if nv != "" {
			l = append(l,
				fmt.Sprintf("%s=%s",
					url.QueryEscape(k),
					url.QueryEscape(nv)))
		} else if k == "acl" {
			l = append(l, url.QueryEscape(k))
		} else if nv == "" {
			l = append(l, url.QueryEscape(k))
		}
	}
	if len(l) != 0 {
		Url = fmt.Sprintf("%s?%s", Url, strings.Join(l, "&"))
	}
	return Url
}

func safeGetElement(name string, contianer map[string]string) string {
	v, ok := contianer[strings.TrimSpace(name)]
	if ok {
		return v
	}
	return ""
}

func getAssign(secret, method string, headers http.Header,
	resource string, result *[]string) string {
	canonicalizedBcHeaders := ""
	canonicalizedResource := resource

	date := headers.Get("Date")
	tmpHeader := formatHeader(headers)
	if len(tmpHeader) > 0 {
		var xHeaderList []string
		for k := range tmpHeader {
			xHeaderList = append(xHeaderList, k)
		}
		sort.Strings(xHeaderList)
		for _, k := range xHeaderList {
			if strings.HasPrefix(k, OasDefineHeaderPrefix) {
				canonicalizedBcHeaders = fmt.Sprintf("%s%s:%v\n",
					canonicalizedBcHeaders, k, safeGetElement(k, tmpHeader))
			}
		}
	}
	stringToSign := fmt.Sprintf("%s\n%s\n%s%s", method, date,
		canonicalizedBcHeaders, canonicalizedResource)
	*result = append(*result, stringToSign)

	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(stringToSign))
	b := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return strings.TrimSpace(b)
}

func formatHeader(headers http.Header) map[string]string {
	tmpHeaders := make(map[string]string)
	for k := range headers {
		kLower := strings.ToLower(k)
		if strings.HasPrefix(kLower, OasDefineHeaderPrefix) {
			tmpHeaders[kLower] = headers.Get(k)
		} else {
			tmpHeaders[k] = headers.Get(k)
		}
	}
	return tmpHeaders
}

func checkResponse(r *http.Response, status int) (err error) {
	if r.StatusCode != status {
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) != 0 {
			reason := new(ErrorMsg)
			err = json.Unmarshal(b, reason)
			if err != nil {
				return
			}
			err = fmt.Errorf(operationFailedFormat, reason)
		} else {
			err = fmt.Errorf("Get HTTP Status: %s", r.Status)
		}
	}
	return
}
