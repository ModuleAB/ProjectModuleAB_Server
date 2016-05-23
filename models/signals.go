package models

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"moduleab_server/common"

	"time"

	"github.com/astaxie/beego"
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

var (
	SignalChannels map[string]chan Signal
)

func init() {
	SignalChannels = make(map[string]chan Signal)
}

type Signal map[string]interface{}

func AddSignal(hostId string, signal Signal) (string, error) {
	keyName := fmt.Sprintf("%s%s", common.DefaultRedisKey, hostId)
	newId := uuid.New()
	signal["id"] = newId
	var (
		buf []byte
		err error
	)
	if !common.DefaultRedisClient.IsExist(keyName) {
		var v = make([]Signal, 0)
		v = append(v, signal)
		// You have 30 minutes to take it out, or failed
		buf, err = toGob(v)
		if err != nil {
			return "", err
		}
	} else {
		b := common.DefaultRedisClient.Get(keyName)
		v, err := fromGob(b.([]byte))
		if err != nil {
			return "", err
		}
		beego.Debug("Got from redis:", v)
		v = append(v, signal)
		buf, err = toGob(v)
		if err != nil {
			return "", err
		}
	}
	return newId, common.DefaultRedisClient.Put(keyName, buf, 30*time.Minute)
}

func GetSignals(hostId string) []Signal {
	keyName := fmt.Sprintf("%s%s", common.DefaultRedisKey, hostId)
	b := common.DefaultRedisClient.Get(keyName)
	beego.Debug("Got from redis:", b)
	v, err := fromGob(b.([]byte))
	if err != nil {
		beego.Warn(err)
		return nil
	}
	return v
}

func GetSignal(hostId, id string) (Signal, error) {
	signals := GetSignals(hostId)
	beego.Debug("Signals", signals)
	for _, v := range signals {
		if v["id"] == id {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Got nothing")
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
		var a = make([]Signal, 0)
		for _, v := range n {
			if v["id"] != signalId {
				a = append(a, v)
			}
		}
		return common.DefaultRedisClient.Put(keyName, a, 30*time.Minute)
	}
	return ErrorSignalNotFound
}

func NotifySignal(hostId, signalId string) error {
	_, ok := SignalChannels[hostId]
	if !ok {
		SignalChannels[hostId] = make(chan Signal, 1024)
	}
	signal, err := GetSignal(hostId, signalId)
	if err != nil {
		return err
	}
	SignalChannels[hostId] <- signal
	return nil
}

func MakeDownloadSignal(path, endpoint, bucket string) Signal {
	s := make(Signal)
	s["type"] = SignalTypeDownload
	s["path"] = path
	s["endpoint"] = endpoint
	s["bucket"] = bucket
	return s
}

func toGob(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(v)
	return buf.Bytes(), err
}

func fromGob(b []byte) ([]Signal, error) {
	var v = make([]Signal, 0)
	var buf = bytes.NewBuffer(b)
	err := gob.NewDecoder(buf).Decode(&v)
	return v, err
}
