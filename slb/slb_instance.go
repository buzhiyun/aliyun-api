package slb

import (
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

var (
	_client = client()
)

func client() *slb.Client {

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
	return _c
}

//func InitSlb() (err error) {
//	if client() == nil {
//		err = errors.New("初始化 slb client 失败")
//	}
//
//	return
//}

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

		response, err := _client.DescribeLoadBalancers(request)

		if err != nil {
			log.Errorf("[slb] 根据ecsId %s 查找 slb失败, %s", ecsServerId, err.Error())
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

	response, err := _client.DescribeLoadBalancerAttribute(request)
	if err != nil {
		log.Errorf("[slb] 获取slb %s 后端服务器失败, %s", slbId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return bkServer, err
	}

	bkServer = append(bkServer, response.BackendServers.BackendServer...)

	return

}

type VServerGroup struct {
	VServerGroupId    string `json:"VServerGroupId"`
	AssociatedObjects struct {
		Listeners struct {
			Listener []struct {
				Port     int    `json:"Port"`
				Protocol string `json:"Protocol"`
			} `json:"Listener"`
		} `json:"Listeners"`
		Rules struct {
			Rule []interface{} `json:"Rule"`
		} `json:"Rules"`
	} `json:"AssociatedObjects"`
	ServiceManagedMode string    `json:"ServiceManagedMode,omitempty"`
	CreateTime         time.Time `json:"CreateTime"`
	VServerGroupName   string    `json:"VServerGroupName"`
	ServerCount        int       `json:"ServerCount"`
}

type DescribeVServerGroupsResponse struct {
	RequestId     string         `json:"RequestId"`
	VServerGroups []VServerGroup `json:"VServerGroups"`
}

// 根据slb去找 虚拟服务器组
func GetSlbVserverGroup(slbId string) (vServerGroups []VServerGroup, err error) {

	request := requests.NewCommonRequest()

	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "slb.aliyuncs.com"
	request.Version = "2014-05-15"
	request.ApiName = "DescribeVServerGroups"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["IncludeRule"] = "true"
	request.QueryParams["IncludeListener"] = "true"
	request.QueryParams["LoadBalancerId"] = slbId

	//request := slb.CreateDescribeVServerGroupsRequest()
	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(60 * time.Second)

	response, err := utils.AliyunClient.ProcessCommonRequest(request)

	if err != nil {
		log.Errorf("[slb] 获取slb %s 的虚拟服务器组失败, %s", slbId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return vServerGroups, err
	}
	var resp DescribeVServerGroupsResponse
	err = json.UnmarshalFromString(response.GetHttpContentString(), &resp)
	if err != nil {
		log.Errorf("[slb] 解析json 返回异常 , %s, %s", err.Error(), response.GetHttpContentString())
	}

	vServerGroups = append(vServerGroups, resp.VServerGroups...)
	return
}

type BackendServers struct {
	Type        string `json:"Type"`
	ServerId    string `json:"ServerId"`
	Port        *int   `json:"Port,omitempty"`
	Weight      int    `json:"Weight"`
	ServerIp    string `json:"ServerIp,omitempty"`
	Description string `json:"Description,omitempty"`
}

type DescribeVServerGroupAttributeResp struct {
	VServerGroupId   string    `json:"VServerGroupId"`
	RequestId        string    `json:"RequestId"`
	CreateTime       time.Time `json:"CreateTime"`
	VServerGroupName string    `json:"VServerGroupName"`
	LoadBalancerId   string    `json:"LoadBalancerId"`
	BackendServers   struct {
		BackendServer []BackendServers `json:"BackendServer"`
	} `json:"BackendServers"`
	Tags struct {
		Tag []interface{} `json:"Tag"`
	} `json:"Tags"`
}

// 根据虚拟服务器组 去找 后端服务器
func GetSlbVserverGroupBackendServer(vServerGroupId string) (bkServer []BackendServers, err error) {

	request := requests.NewCommonRequest()
	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(45 * time.Second)

	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "slb.aliyuncs.com"
	request.Version = "2014-05-15"
	request.ApiName = "DescribeVServerGroupAttribute"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["VServerGroupId"] = vServerGroupId

	response, err := utils.AliyunClient.ProcessCommonRequest(request)
	if err != nil {
		log.Errorf("[slb] 获取slb虚拟服务器 %s 组详情失败, %s", vServerGroupId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return bkServer, err
	}

	var resp DescribeVServerGroupAttributeResp
	err = json.UnmarshalFromString(response.GetHttpContentString(), &resp)
	if err != nil {
		log.Errorf("[slb] 解析json 返回异常 , %s, %s", err.Error(), response.GetHttpContentString())
	}

	bkServer = append(bkServer, resp.BackendServers.BackendServer...)
	return
}

//type SetBackendServersResp struct {
//	RequestId      string `json:"RequestId"`
//	LoadBalancerId string `json:"LoadBalancerId"`
//	BackendServers struct {
//		BackendServer []backendServer `json:"BackendServer"`
//	} `json:"BackendServers"`
//}

// 设置SLB后端服务器权重
func SetSlbBackendServer(slbId string, backendServers []backendServer) (err error) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "slb.aliyuncs.com"
	request.Version = "2014-05-15"
	request.ApiName = "SetBackendServers"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["LoadBalancerId"] = slbId

	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(60 * time.Second)

	bkserverJson, err := json.MarshalToString(backendServers)
	if err != nil {
		log.Errorf("[slb] 解析backendServers数据失败 %#v 权重失败, %s", backendServers, err.Error())
		return
	}
	log.Debugf("bkserverJson: %s", bkserverJson)
	request.QueryParams["BackendServers"] = bkserverJson

	response, err := utils.AliyunClient.ProcessCommonRequest(request)
	if err != nil {
		log.Errorf("[slb] 设置权重 %s 权重失败, %s", slbId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return
	}

	log.Infof("设置 %s 权重 %s , %s", slbId, backendServers, response.GetHttpContentString())

	return
}

// 这个是传参过去的对象，slb 的接口组装的时候  port 和weight 在json里必须是string
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
	request := requests.NewCommonRequest()
	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(60 * time.Second)

	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "slb.aliyuncs.com"
	request.Version = "2014-05-15"
	request.ApiName = "SetVServerGroupAttribute"
	request.QueryParams["RegionId"] = "cn-hangzhou"

	bkserverJson, err := json.MarshalToString(backendServers)
	log.Debugf("bkserverJson: %s", bkserverJson)

	if err != nil {
		log.Errorf("[slb] 解析backendServers数据失败 %#v 权重失败, %s", backendServers, err.Error())
		return
	}
	request.QueryParams["BackendServers"] = bkserverJson
	request.QueryParams["VServerGroupId"] = vGroupId

	response, err := utils.AliyunClient.ProcessCommonRequest(request)

	if err != nil {
		log.Errorf("[slb] 设置虚拟服务器组权重 %s 权重失败, %s", vGroupId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return
	}

	log.Infof("[slb] 设置虚拟服务器组 %s 权重 %s \n%s", vGroupId, backendServers, response.GetHttpContentString())
	return
}
