package ecs

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/buzhiyun/aliyun-api/msg"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/buzhiyun/go-utils/log"
	"time"
)

var (
	ecsInstances *[]ecs.Instance
	_client      *ecs.Client
)

func init() {
	ecsInstances = &[]ecs.Instance{}
}

func client() *ecs.Client {
	if _client != nil {
		return _client
	}

	regionId, aliyunKey, aliyunSecret, _ := utils.GetAliyunKey()

	config := sdk.NewConfig()
	// 是否开启重试机制
	config.WithAutoRetry(true)
	// 最大重试次数
	config.WithMaxRetryTime(3)

	credential := credentials.NewAccessKeyCredential(aliyunKey, aliyunSecret)
	_c, err := ecs.NewClientWithOptions(regionId, config, credential)

	if err != nil {
		log.Errorf("初始化 ecs client 失败, %s", err.Error())
		return nil
	}
	_client = _c
	return _c
}

func InitECS() (err error) {
	if client() == nil {
		err = errors.New("初始化 ecs client 失败")
	}

	return
}

// 得到所有ECS 主机
func GetInstances() (res []ecs.Instance, err error) {
	var instances []ecs.Instance

	pageNum := 1 // 先查第一页的
	maxPage := 1 //默认最大页数就是1
	for pageNum <= maxPage {
		request := ecs.CreateDescribeInstancesRequest()
		request.PageSize = requests.NewInteger(100)
		request.PageNumber = requests.NewInteger(pageNum)
		request.SetReadTimeout(60 * time.Second)
		request.SetConnectTimeout(5 * time.Second)
		response, err := client().DescribeInstances(request)

		if err != nil {
			log.Error(err.Error())
			msg.AliyunSdkAlert(err.Error())
			return instances, err
		}

		if response != nil {
			maxPage = ((response.TotalCount - 1) / 100) + 1

			instances = append(instances, response.Instances.Instance...)
		}
		//增加页码，准备取下一页
		pageNum++
	}

	return instances, nil
}

// 返回ecs的所有IP信息
func GetInstanceIp(instance ecs.Instance) (ipAddresses []string) {

	ipAddresses = append(ipAddresses, instance.InnerIpAddress.IpAddress...)
	ipAddresses = append(ipAddresses, instance.PublicIpAddress.IpAddress...)
	ipAddresses = append(ipAddresses, instance.VpcAttributes.PrivateIpAddress.IpAddress...)

	return ipAddresses
}

// 返回ECS 的私网IP

func GetInstancesPrivateIP(instance ecs.Instance) string {
	var ipAddresses []string
	ipAddresses = append(ipAddresses, instance.InnerIpAddress.IpAddress...)
	ipAddresses = append(ipAddresses, instance.VpcAttributes.PrivateIpAddress.IpAddress...)

	if len(ipAddresses) == 0 {
		panic(instance.InstanceName + " 没有私网IP")
	}

	return ipAddresses[0]
}

/*
刷新ECS方法
*/
func UpdateEcs() (refreshCount int, err error) {

	instances, err := GetInstances()
	if err != nil {
		log.Errorf("刷新异常 %s", err.Error())
		return
	}

	ecsInstances = &instances

	refreshCount = len(instances)
	log.Infof("刷新了 %v 条记录", refreshCount)
	//logger.Println("刷新了" + strconv.Itoa(count) + "条记录" )
	return refreshCount, err
}

// 按主机名搜索
func SearchByName(hostname string) (res []ecs.Instance) {
	for _, instance := range *ecsInstances {
		if utils.MatchWildcard(instance.InstanceName, hostname) {
			res = append(res, instance)
		}
	}
	return
}

// 按IP搜索
func SearchByIP(searchIP string) (res []ecs.Instance) {
	for _, instance := range *ecsInstances {
		for _, ip := range GetInstanceIp(instance) {
			if utils.MatchWildcard(ip, searchIP) {
				res = append(res, instance)
				break
			}
		}
	}
	return
}
