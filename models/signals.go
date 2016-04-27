package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
)

const (
	defaultRedisKey   = "ModuleAB"
	CacheSignalPrefix = "Signal_"
)

const (
	SignalTypeNothing = iota
	SignalTypeDownload
)

var redis cache.Cache

func init() {
	var err error
	redisConf := make(map[string]string)
	redisConf["conn"] = beego.AppConfig.String("redis::host")
	redisConf["password"] = beego.AppConfig.String("reids::password")
	redisConf["key"] = beego.AppConfig.String("redis::key")
	b, _ := json.Marshal(redisConf)
	redis, err = cache.NewCache("redis", string(b))
	if err != nil {
		beego.Alert("Connect to redis failed:", err)
	}
}

type Signal map[string]interface{}

func AddSignal(hostId string, signal Signal) error {
	keyName := fmt.Sprintf("%s%s", defaultRedisKey, hostId)
	if !redis.IsExist(keyName) {
		v := make([]Signal, 0)
		v = append(v, signal)
		// You have 30 minutes to take it out, or failed
		return redis.Put(keyName, v, 30*time.Minute)

	} else {
		v := redis.Get(keyName)
		n, ok := v.([]Signal)
		if !ok {
			return fmt.Errorf("Bad DataType")
		}
		n = append(n, signal)
		return redis.Put(keyName, v, 30*time.Minute)
	}
	return nil
}

func GetSignals(hostId string) []Signal {
	keyName := fmt.Sprintf("%s%s", defaultRedisKey, hostId)
	v := redis.Get(keyName)
	n, ok := v.([]Signal)
	if !ok {
		return nil
	}
	return n
}

func TruncateSignals(hostId string) {
	keyName := fmt.Sprintf("%s%s", defaultRedisKey, hostId)
	redis.Delete(keyName)
}

func DeleteSignal(hostId string, signal Signal) error {
	keyName := fmt.Sprintf("%s%s", defaultRedisKey, hostId)
	if redis.IsExist(keyName) {
		v := redis.Get(keyName)
		n, ok := v.([]Signal)
		if !ok {
			return fmt.Errorf("Bad DataType")
		}
		a := make([]Signal, 0)
		for _, v := range n {
			if equal(v, signal) {
				a = append(a, v)
			}
		}
		return redis.Put(keyName, a, 30*time.Minute)
	}
	return nil
}

func equal(x, y Signal) bool {
	if len(x) != len(y) {
		return false
	}
	for k, xv := range x {
		if yv, ok := y[k]; !ok || yv != xv {
			return false
		}
	}
	return true
}
