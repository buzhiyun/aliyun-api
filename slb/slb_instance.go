package slb

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/buzhiyun/aliyun-api/msg"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/buzhiyun/go-utils/log"
	"time"
)

//type backendServer struct {
//	Type     string `json:"Type"`
//	ServerId string `json:"ServerId"`
//	Port     *int    `json:"Port,omitempty"`
//	Weight   int    `json:"Weight"`
//}

var _client *slb.Client

func client() *slb.Client {
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
	_c, err := slb.NewClientWithOptions(regionId, config, credential)

	if err != nil {
		log.Errorf("初始化 slb client 失败, %s", err.Error())
		return nil
	}
	_client = _c
	return _c
}

func InitSlb() (err error) {
	if client() == nil {
		err = errors.New("初始化 slb client 失败")
	}

	return
}

// 查找所有有该ECS的slb
func GetEcsSlb(ecsServerId string) (slbs []slb.LoadBalancer, err error) {

	if err != nil {
		log.Errorf("初始化 slb client 失败, %s", err.Error())
		return
	}

	pageNum := 1 // 先查第一页的
	maxPage := 1 //默认最大页数就是1

	for pageNum <= maxPage {
		request := slb.CreateDescribeLoadBalancersRequest()
		// 连接超时设置，仅对当前请求有效。
		request.SetConnectTimeout(5 * time.Second)
		// 读超时设置，仅对当前请求有效。
		request.SetReadTimeout(60 * time.Second)

		request.Scheme = "https"

		request.ServerId = ecsServerId
		request.PageSize = requests.NewInteger(100)

		response, err := client().DescribeLoadBalancers(request)

		if err != nil {
			log.Errorf("根据ecsId %s 查找 slb失败, %s", ecsServerId, err.Error())
			msg.AliyunSdkAlert(err.Error())
			return slbs, err
		}

		if response != nil {
			maxPage = ((response.TotalCount - 1) / 100) + 1

			slbs = append(slbs, response.LoadBalancers.LoadBalancer...)
		}
		//增加页码，准备取下一页
		pageNum++
	}

	return

}

// 根据slb去找 后端服务器 【不是虚拟服务器组】
func GetSlbBackendServer(slbId string) (bkServer []slb.BackendServerInDescribeLoadBalancerAttribute, err error) {

	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(60 * time.Second)

	request.Scheme = "https"

	request.LoadBalancerId = slbId

	response, err := client().DescribeLoadBalancerAttribute(request)
	if err != nil {
		log.Errorf("获取slb %s 后端服务器失败, %s", slbId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return bkServer, err
	}

	bkServer = append(bkServer, response.BackendServers.BackendServer...)

	return

}

// 根据slb去找 虚拟服务器组
func GetSlbVserverGroup(slbId string) (vServerGroups []slb.VServerGroup, err error) {

	request := slb.CreateDescribeVServerGroupsRequest()
	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(60 * time.Second)
	request.Scheme = "https"

	request.LoadBalancerId = slbId

	response, err := client().DescribeVServerGroups(request)
	if err != nil {
		log.Errorf("获取slb %s 的虚拟服务器组失败, %s", slbId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return vServerGroups, err
	}

	vServerGroups = append(vServerGroups, response.VServerGroups.VServerGroup...)
	return
}

// 根据虚拟服务器组 去找 后端服务器
func GetSlbVserverGroupBackendServer(vServerGroupId string) (bkServer []slb.BackendServerInDescribeVServerGroupAttribute, err error) {

	request := slb.CreateDescribeVServerGroupAttributeRequest()
	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(60 * time.Second)

	request.Scheme = "https"

	request.VServerGroupId = vServerGroupId

	response, err := client().DescribeVServerGroupAttribute(request)
	if err != nil {
		log.Errorf("获取slb虚拟服务器 %s 组详情失败, %s", vServerGroupId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return bkServer, err
	}

	bkServer = append(bkServer, response.BackendServers.BackendServer...)
	return
}

// 设置后端服务器
func SetSlbBackendServer(slbId string, backendServers []backendServer) (err error) {

	request := slb.CreateSetBackendServersRequest()
	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(60 * time.Second)
	request.Scheme = "https"

	request.LoadBalancerId = slbId

	bkserverJson, err := json.MarshalToString(backendServers)
	if err != nil {
		log.Errorf("解析backendServers数据失败 %#v 权重失败, %s", backendServers, err.Error())
		return
	}
	log.Debugf("bkserverJson: %s", bkserverJson)
	request.BackendServers = bkserverJson

	response, err := client().SetBackendServers(request)
	if err != nil {
		log.Errorf("设置权重 %s 权重失败, %s", slbId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return
	}

	log.Infof("设置 %s 权重 %s , %s", slbId, backendServers, response.GetHttpContentString())

	return
}

type backendServer struct {
	ServerId    string `json:"ServerId"`
	Weight      string `json:"Weight"`
	Type        string `json:"Type"`
	ServerIp    string `json:"ServerIp,omitempty"`
	Port        string `json:"Port,omitempty"`
	Description string `json:"Description,omitempty"`
}

// 设置后端虚拟服务器组
func SetSlbVserverGroup(vGroupId string, backendServers []backendServer) (err error) {
	request := slb.CreateSetVServerGroupAttributeRequest()
	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(60 * time.Second)

	request.Scheme = "https"

	bkserverJson, err := json.MarshalToString(backendServers)
	log.Debugf("bkserverJson: %s", bkserverJson)

	if err != nil {
		log.Errorf("解析backendServers数据失败 %#v 权重失败, %s", backendServers, err.Error())
		return
	}
	request.BackendServers = bkserverJson
	request.VServerGroupId = vGroupId

	response, err := client().SetVServerGroupAttribute(request)
	if err != nil {
		log.Errorf("设置虚拟服务器组权重 %s 权重失败, %s", vGroupId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return
	}

	log.Infof("设置虚拟服务器组 %s 权重 %s \n%s", vGroupId, backendServers, response.GetHttpContentString())
	return
}
