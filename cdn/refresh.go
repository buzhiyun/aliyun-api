package cdn

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/kataras/golog"
	"strings"
)


var _client *cdn.Client

func client() *cdn.Client {
	if _client != nil {
		return _client
	}

	regionId, aliyunKey, aliyunSecret ,_ := utils.GetAliyunKey()

	config := sdk.NewConfig()
	// 是否开启重试机制
	config.WithAutoRetry(true);
	// 最大重试次数
	config.WithMaxRetryTime(3);

	credential := credentials.NewAccessKeyCredential(aliyunKey, aliyunSecret)
	_c, err := cdn.NewClientWithOptions(regionId, config, credential)

	if err != nil {
		golog.Errorf("初始化 cdn client 失败, %s",err.Error())
		return nil
	}
	_client = _c
	return _c
}


func InitCDN() (err error)  {
	if client() == nil {
		err = errors.New("初始化 cdn client 失败")
	}

	return
}



// 刷新CDN缓存
func RefreshUrl(urls []string) (response *cdn.RefreshObjectCachesResponse,err error) {


	request := cdn.CreateRefreshObjectCachesRequest()

	request.Scheme = "https"

	request.ObjectPath = strings.Join(urls,"\n")

	response, err = client().RefreshObjectCaches(request)
	if err != nil {
		golog.Errorf("刷新cdn失败 %s",err.Error())
	}
	golog.Infof("刷新cdn成功 %s", response.GetHttpContentString())

	return
}


// 预热cdn
func PushObjectCache(urls []string) (response *cdn.PushObjectCacheResponse, err error) {

	request := cdn.CreatePushObjectCacheRequest()

	request.Scheme = "https"

	request.Area = "domestic"
	request.L2Preload = requests.NewBoolean(true)
	request.ObjectPath = strings.Join(urls,"\n")


	response, err = client().PushObjectCache(request)
	if err != nil {
		golog.Errorf("预热cdn失败 %s",err.Error())
	}

	golog.Infof("预热cdn成功 %s", response.GetHttpContentString())

	return

}