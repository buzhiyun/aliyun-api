package cms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/buzhiyun/go-utils/log"
)

type aliyunCms struct {
	client *cms.Client
}

var (
	CMS = aliyunCms{
		client: newclient(),
	}
)

func newclient() *cms.Client {

	regionId, aliyunKey, aliyunSecret, _ := utils.GetAliyunKey()

	config := sdk.NewConfig()
	// 是否开启重试机制
	config.WithAutoRetry(true)
	// 最大重试次数
	config.WithMaxRetryTime(3)

	credential := credentials.NewAccessKeyCredential(aliyunKey, aliyunSecret)
	_c, err := cms.NewClientWithOptions(regionId, config, credential)
	if err != nil {
		log.Fatalf("初始化 cms client 失败, %s", err.Error())
		return nil
	}
	return _c

}
