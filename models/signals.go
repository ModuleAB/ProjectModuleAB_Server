package models

import (
	"errors"
	"fmt"
	"moduleab_server/common"
	"time"

	"github.com/pborman/uuid"
)

const (
	CacheSignalPrefix = "Signal_"
)

const (
	SignalTypeNothing = iota
	SignalTypeDownload
)

var (
	ErrorSignalNotFound    = errors.New("Signal Not Found")
	ErrorSignalBadDataType = errors.New("Bad data type")
)

type Signal map[string]interface{}

func AddSignal(hostId string, signal Signal) (string, error) {
	keyName := fmt.Sprintf("%s%s", common.DefaultRedisKey, hostId)
	newId := uuid.New()
	signal["id"] = newId
	if !common.DefaultRedisClient.IsExist(keyName) {
		v := make([]Signal, 0)
		v = append(v, signal)
		// You have 30 minutes to take it out, or failed
		return newId, common.DefaultRedisClient.Put(keyName, v, 30*time.Minute)

	} else {
		v := common.DefaultRedisClient.Get(keyName)
		n, ok := v.([]Signal)
		if !ok {
			return "", ErrorSignalBadDataType
		}
		n = append(n, signal)
		return newId, common.DefaultRedisClient.Put(keyName, v, 30*time.Minute)
	}
	return "", nil
}

func GetSignals(hostId string) []Signal {
	keyName := fmt.Sprintf("%s%s", common.DefaultRedisKey, hostId)
	v := common.DefaultRedisClient.Get(keyName)
	n, ok := v.([]Signal)
	if !ok {
		return nil
	}
	return n
}

func TruncateSignals(hostId string) {
	keyName := fmt.Sprintf("%s%s", common.DefaultRedisKey, hostId)
	common.DefaultRedisClient.Delete(keyName)
}

func DeleteSignal(hostId string, signalId string) error {
	keyName := fmt.Sprintf("%s%s", common.DefaultRedisKey, hostId)
	if common.DefaultRedisClient.IsExist(keyName) {
		v := common.DefaultRedisClient.Get(keyName)
		n, ok := v.([]Signal)
		if !ok {
			return ErrorSignalBadDataType
		}
		a := make([]Signal, 0)
		for _, v := range n {
			if v["id"] != signalId {
				a = append(a, v)
			}
		}
		return common.DefaultRedisClient.Put(keyName, a, 30*time.Minute)
	}
	return ErrorSignalNotFound
}
