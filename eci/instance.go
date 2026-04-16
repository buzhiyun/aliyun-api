package eci

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/buzhiyun/go-utils/log"
)

var (
	Client = NewEciClient()
)

type eciClient struct {
	client       *eci.Client
	ecsInstances []eci.DescribeContainerGroupsContainerGroup0
}

func NewEciClient() *eciClient {
	cli := &eciClient{
		client: client(),
	}

	go func(c *eciClient) {
		for {
			cli.RefreshEciContainer()
			time.Sleep(180 * time.Second)
		}
	}(cli)

	return cli
}

func (c *eciClient) RefreshEciContainer() (refreshCount int, err error) {
	c.ecsInstances = nil
	request := eci.CreateDescribeContainerGroupsRequest()
	request.Scheme = "https"
	request.ContainerGroupName = ""
	request.VSwitchId = ""
	request.Status = ""

	response, err := c.client.DescribeContainerGroups(request)
	if err != nil {
		log.Errorf("[eci] 刷新容器组列表失败, %s", err.Error())
		return
	}

	c.ecsInstances = response.ContainerGroups
	// 懒加载查询
	for len(response.NextToken) > 0 {

		request.NextToken = response.NextToken
		response, err = c.client.DescribeContainerGroups(request)
		if err != nil {
			log.Errorf("[eci] 刷新容器组列表失败, %s", err.Error())
			return
		}
		c.ecsInstances = append(c.ecsInstances, response.ContainerGroups...)
	}
	refreshCount = len(c.ecsInstances)
	log.Infof("[eci] 刷新容器组列表成功, 共 %d 个实例", refreshCount)
	return
}

func client() *eci.Client {
	regionId, aliyunKey, aliyunSecret, _ := utils.GetAliyunKey()

	config := sdk.NewConfig()
	// 是否开启重试机制
	config.WithAutoRetry(true)
	// 最大重试次数
	config.WithMaxRetryTime(3)

	credential := credentials.NewAccessKeyCredential(aliyunKey, aliyunSecret)
	_c, err := eci.NewClientWithOptions(regionId, config, credential)

	if err != nil {
		log.Errorf("初始化 ecs client 失败, %s", err.Error())
		return nil
	}
	return _c
}

type searchType string

const (
	SearchName searchType = "name"
	SearchIP   searchType = "ip"
)

func (c *eciClient) SearchContainer(keyword string, searchType searchType) (ContainerGroups []eci.DescribeContainerGroupsContainerGroup0) {
	if keyword == "" {
		return c.ecsInstances
	}

	for _, containerGroup := range c.ecsInstances {
		if (searchType == SearchName && utils.MatchWildcard(containerGroup.ContainerGroupName, keyword)) ||
			(searchType == SearchIP && utils.MatchWildcard(containerGroup.IntranetIp, keyword)) {
			ContainerGroups = append(ContainerGroups, containerGroup)
		}
	}
	return
}
