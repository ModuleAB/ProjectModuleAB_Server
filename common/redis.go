package common

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
)

const DefaultRedisKey = "ModuleAB"

var DefaultRedisClient cache.Cache

func init() {
	var err error
	redisConf := make(map[string]string)
	redisConf["conn"] = beego.AppConfig.String("redis::host")
	redisConf["password"] = beego.AppConfig.String("reids::password")
	redisConf["key"] = beego.AppConfig.String("redis::key")
	b, _ := json.Marshal(redisConf)
	DefaultRedisClient, err = cache.NewCache("redis", string(b))
	if err != nil {
		beego.Alert("Connect to redis failed:", err)
	}
}
