package utils

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/buzhiyun/go-utils/cfg"
	"github.com/buzhiyun/go-utils/log"
)

// 获取阿里云的 key secret
func GetAliyunKey() (regionId, aliyunKey, aliyunSecret string, err error) {
	regionId, ok := cfg.Config().GetString("aliyun.region")
	if !ok {
		err = errors.New("读取配置 aliyun.region 异常")
		log.Error(err.Error())
		regionId = "cn-hangzhou" // 默认使用cn-hangzhou
	}

	aliyunKey, ok = cfg.Config().GetString("aliyun.key")
	if !ok {
		err = errors.New("读取配置 aliyun.key 异常")
		log.Fatal(err.Error())
		return
	}

	aliyunSecret, ok = cfg.Config().GetString("aliyun.secret")
	if !ok {
		err = errors.New("读取配置 aliyun.secret 异常")
		log.Fatal(err.Error())
		return
	}
	return
}

var AliyunClient = client()

func client() *sdk.Client {
	regionId, aliyunKey, aliyunSecret, _ := GetAliyunKey()

	config := sdk.NewConfig()
	// 是否开启重试机制
	config.WithAutoRetry(true)
	// 最大重试次数
	config.WithMaxRetryTime(3)

	credential := credentials.NewAccessKeyCredential(aliyunKey, aliyunSecret)
	_c, err := sdk.NewClientWithOptions(regionId, config, credential)
	if err != nil {
		log.Errorf("初始化 slb client 失败, %s", err.Error())
		return nil
	}
	return _c
}
