package docs

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/swagger"
)

const (
    Rootinfo string = `{"apiVersion":"1.0.0","swaggerVersion":"1.2","apis":[{"path":"/hosts","description":""},{"path":"/client","description":""},{"path":"/appSets","description":""},{"path":"/oss","description":""},{"path":"/oas","description":""}],"info":{"title":"ModuleAB API","description":"ModuleAB server API","contact":"tonychyi1989@gmail.com","license":"Url http://www.gnu.org/licenses/gpl-3.0.html"}}`
    Subapi string = `{"/appSets":{"apiVersion":"1.0.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/appSets","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"createAppSet","type":""}]},{"path":"/:name","description":"","operations":[{"httpMethod":"GET","nickname":"getAppSet","type":""}]},{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"listAppSets","type":""}]},{"path":"/:name","description":"","operations":[{"httpMethod":"DELETE","nickname":"deleteAppSet","type":""}]},{"path":"/:name","description":"","operations":[{"httpMethod":"PUT","nickname":"updateAppSet","type":""}]}]},"/client":{"apiVersion":"1.0.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/client","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/config","description":"","operations":[{"httpMethod":"GET","nickname":"getClientConf","type":"","summary":"getClientConf","parameters":[{"paramType":"body","name":"body","description":"\"body for host content\"","dataType":"Hosts","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":0,"message":"200","responseModel":""},{"code":403,"message":"body is empty","responseModel":""}]}]},{"path":"/signal/:name","description":"","operations":[{"httpMethod":"GET","nickname":"getSignals","type":""}]},{"path":"/signal/:name/:id","description":"","operations":[{"httpMethod":"GET","nickname":"getSignal","type":""}]},{"path":"/signal/:name","description":"","operations":[{"httpMethod":"POST","nickname":"createSignal","type":""}]},{"path":"/signal/:name/:id","description":"","operations":[{"httpMethod":"DELETE","nickname":"deleteSignal","type":""}]}]},"/hosts":{"apiVersion":"1.0.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/hosts","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"createHost","type":"","summary":"create Host","parameters":[{"paramType":"body","name":"host","description":"\"host\"","dataType":"object","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":201,"message":"models.Hosts","responseModel":"Hosts"},{"code":400,"message":"Hostname or IP missing","responseModel":""},{"code":500,"message":"Failure on writing database","responseModel":""}]}]},{"path":"/:name","description":"","operations":[{"httpMethod":"GET","nickname":"getHost","type":"","summary":"get Host info","parameters":[{"paramType":"body","name":"body","description":"\"body for host content\"","dataType":"Hosts","type":"","format":"","allowMultiple":false,"required":true,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"models.Hosts.Id","responseModel":""},{"code":403,"message":"body is empty","responseModel":""}]}]},{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"listHosts","type":"","summary":"get all Host info","responseMessages":[{"code":0,"message":"200","responseModel":""}]}]},{"path":"/:name","description":"","operations":[{"httpMethod":"DELETE","nickname":"deleteHost","type":"","summary":"delete host","responseMessages":[{"code":0,"message":"204","responseModel":""},{"code":404,"message":"","responseModel":""}]}]},{"path":"/:name","description":"","operations":[{"httpMethod":"PUT","nickname":"updateHost","type":"","summary":"update host","responseMessages":[{"code":0,"message":"204","responseModel":""},{"code":404,"message":"","responseModel":""}]}]}],"models":{"AppSets":{"id":"AppSets","properties":{"Hosts":{"type":"array","description":"","items":{"$ref":"Hosts"},"format":""},"Policies":{"type":"array","description":"","items":{"$ref":"Policies"},"format":""},"Records":{"type":"array","description":"","items":{"$ref":"Records"},"format":""},"description":{"type":"string","description":"","format":""},"id":{"type":"string","description":"","format":""},"name":{"type":"string","description":"","format":""}}},"BackupSets":{"id":"BackupSets","properties":{"Desc":{"type":"string","description":"","format":""},"Hosts":{"type":"array","description":"","items":{"$ref":"Hosts"},"format":""},"Id":{"type":"string","description":"","format":""},"Name":{"type":"string","description":"","format":""},"Oas":{"type":"Oas","description":"","format":""},"Oss":{"type":"Oss","description":"","format":""},"Policies":{"type":"array","description":"","items":{"$ref":"Policies"},"format":""}}},"Hosts":{"id":"Hosts","properties":{"app_set":{"type":"AppSets","description":"","format":""},"backup_sets":{"type":"array","description":"","items":{"$ref":"BackupSets"},"format":""},"id":{"type":"string","description":"","format":""},"ip":{"type":"string","description":"","format":""},"name":{"type":"string","description":"","format":""}}},"Oas":{"id":"Oas","properties":{"BackupSets":{"type":"array","description":"","items":{"$ref":"BackupSets"},"format":""},"Jobs":{"type":"array","description":"","items":{"$ref":"OasJobs"},"format":""},"VaultId":{"type":"string","description":"","format":""},"VaultName":{"type":"string","description":"","format":""},"endpoint":{"type":"string","description":"","format":""},"id":{"type":"string","description":"","format":""}}},"OasJobs":{"id":"OasJobs","properties":{"Id":{"type":"string","description":"","format":""},"JobId":{"type":"string","description":"","format":""},"JobType":{"type":"int","description":"","format":""},"Records":{"type":"Records","description":"","format":""},"RequestId":{"type":"string","description":"","format":""},"Vault":{"type":"Oas","description":"","format":""}}},"Oss":{"id":"Oss","properties":{"BackupSets":{"type":"array","description":"","items":{"$ref":"BackupSets"},"format":""},"id":{"type":"string","description":"","format":""},"name":{"type":"string","description":"","format":""},"region":{"type":"string","description":"","format":""}}},"Policies":{"id":"Policies","properties":{"Action":{"type":"int","description":"","format":""},"AppSet":{"type":"AppSets","description":"","format":""},"BackupSet":{"type":"BackupSets","description":"","format":""},"Desc":{"type":"string","description":"","format":""},"Id":{"type":"string","description":"","format":""},"Name":{"type":"string","description":"","format":""},"ReservePeriod":{"type":"int","description":"","format":""},"Target":{"type":"int","description":"","format":""},"TargetEnd":{"type":"int","description":"","format":""},"TargetStart":{"type":"int","description":"","format":""}}},"Records":{"id":"Records","properties":{"AppSet":{"type":"AppSets","description":"","format":""},"ArchiveId":{"type":"string","description":"","format":""},"ArchivedTime":{"type":"\u0026{time Time}","description":"","format":""},"BackupSet":{"type":"BackupSets","description":"","format":""},"BackupTime":{"type":"\u0026{time Time}","description":"","format":""},"Host":{"type":"Hosts","description":"","format":""},"Id":{"type":"string","description":"","format":""},"Jobs":{"type":"array","description":"","items":{"$ref":"OasJobs"},"format":""},"Path":{"type":"string","description":"","format":""},"Type":{"type":"int","description":"","format":""}}}}},"/oas":{"apiVersion":"1.0.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/oas","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"createOAS","type":""}]},{"path":"/:name","description":"","operations":[{"httpMethod":"GET","nickname":"getOAS","type":""}]},{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"listOAS","type":""}]},{"path":"/:name","description":"","operations":[{"httpMethod":"DELETE","nickname":"deleteOAS","type":""}]},{"path":"/:name","description":"","operations":[{"httpMethod":"PUT","nickname":"updateOAS","type":""}]}]},"/oss":{"apiVersion":"1.0.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/oss","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/","description":"","operations":[{"httpMethod":"POST","nickname":"createOSS","type":""}]},{"path":"/:name","description":"","operations":[{"httpMethod":"GET","nickname":"getOSS","type":""}]},{"path":"/","description":"","operations":[{"httpMethod":"GET","nickname":"listOSS","type":""}]},{"path":"/:name","description":"","operations":[{"httpMethod":"DELETE","nickname":"deleteOSS","type":""}]},{"path":"/:name","description":"","operations":[{"httpMethod":"PUT","nickname":"updateOSS","type":""}]}]}}`
    BasePath string= "/api/v1"
)

var rootapi swagger.ResourceListing
var apilist map[string]*swagger.APIDeclaration

func init() {
	if beego.BConfig.WebConfig.EnableDocs {
		err := json.Unmarshal([]byte(Rootinfo), &rootapi)
		if err != nil {
			beego.Error(err)
		}
		err = json.Unmarshal([]byte(Subapi), &apilist)
		if err != nil {
			beego.Error(err)
		}
		beego.GlobalDocAPI["Root"] = rootapi
		for k, v := range apilist {
			for i, a := range v.APIs {
				a.Path = urlReplace(k + a.Path)
				v.APIs[i] = a
			}
			v.BasePath = BasePath
			beego.GlobalDocAPI[strings.Trim(k, "/")] = v
		}
	}
}


func urlReplace(src string) string {
	pt := strings.Split(src, "/")
	for i, p := range pt {
		if len(p) > 0 {
			if p[0] == ':' {
				pt[i] = "{" + p[1:] + "}"
			} else if p[0] == '?' && p[1] == ':' {
				pt[i] = "{" + p[2:] + "}"
			}
		}
	}
	return strings.Join(pt, "/")
}
