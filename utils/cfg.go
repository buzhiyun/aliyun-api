package utils

import (
	"errors"
	"github.com/buzhiyun/go-utils/cfg"
	"github.com/kataras/golog"
)


// 获取阿里云的 key secret
func GetAliyunKey() (regionId,aliyunKey,aliyunSecret string, err error) {
	regionId , ok := cfg.Config().GetString("aliyun.region")
	if !ok {
		err = errors.New("读取配置 aliyun.region 异常")
		golog.Error(err.Error())
		regionId = "cn-hangzhou"  // 默认使用cn-hangzhou
	}

	aliyunKey , ok = cfg.Config().GetString("aliyun.key")
	if !ok {
		err = errors.New("读取配置 aliyun.key 异常")
		golog.Fatal(err.Error())
		return
	}

	aliyunSecret , ok = cfg.Config().GetString("aliyun.secret")
	if !ok {
		err = errors.New("读取配置 aliyun.secret 异常")
		golog.Fatal(err.Error())
		return
	}
	return
}