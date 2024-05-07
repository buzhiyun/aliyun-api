package slb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/buzhiyun/go-utils/cfg"
	"github.com/buzhiyun/go-utils/log"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

const (
	BackServer = iota
	VgroupBackServer
)

var (
	json        = jsoniter.ConfigCompatibleWithStandardLibrary
	vgroupWhite = getVgroupWhite() // vgroup 白名单
)

func getVgroupWhite() map[string]bool {
	vgWhiteList, ok := cfg.Config().GetStrings("slb.whitelist.vgroup")
	if !ok {
		log.Fatal("[slb] 获取配置 slb.whitelist.vgroup 异常")
		return nil
	}
	var whiteMap = make(map[string]bool)
	for _, s := range vgWhiteList {
		whiteMap[s] = true
		log.Infof("[slb] 添加 vgroup %s 到白名单排除权重设置", s)
	}
	return whiteMap
}

type EcsSetResult struct {
	ServerId  string
	SlbId     string
	SlbName   string
	GroupType int // vgroup
	GroupName string
	From      int
	To        int
}

type setWeightMsg struct {
	Result []EcsSetResult
	Err    error
}

func updateEcsSlbWeight(serverId string, slb slb.LoadBalancer, weight int, resultChan *chan setWeightMsg) {
	var msg setWeightMsg
	defer func(c *chan setWeightMsg, m *setWeightMsg) { *c <- *m }(resultChan, &msg)

	// 查询是否在slb的后端服务器里
	bkServers, err := GetSlbBackendServer(slb.LoadBalancerId)
	if err != nil {
		msg.Err = err
		return
	}
	for _, server := range bkServers {
		if server.ServerId == serverId {

			// 后端服务器设置权重
			newSet := make([]backendServer, len(bkServers))
			// 设置 vServerGroup 里的服务器权重

			for i, _bkServer := range bkServers {
				newSet[i].ServerId = _bkServer.ServerId
				//newSet[i].ServerIp = _bkServer.ServerIp
				//newSet[i].Type = _bkServer.Type
				//newSet[i].Description = _bkServer.Description
				newSet[i].Weight = strconv.Itoa(_bkServer.Weight)
				if _bkServer.ServerId == serverId {
					newSet[i].Weight = strconv.Itoa(weight)
				}
			}

			if err = SetSlbBackendServer(slb.LoadBalancerId, newSet); err != nil {
				continue
			}

			msg.Result = append(msg.Result, EcsSetResult{
				ServerId:  serverId,
				SlbId:     slb.LoadBalancerId,
				SlbName:   slb.LoadBalancerName,
				GroupType: BackServer,
				GroupName: "",
				From:      server.Weight,
				To:        weight,
			})
			break
		}
	}

	// 检查 vgroup 组里的机器
	vGroups, err := GetSlbVserverGroup(slb.LoadBalancerId)
	if err != nil {
		msg.Err = err
		return
	}

	for _, vGroup := range vGroups {
		if _, ok := vgroupWhite[vGroup.VServerGroupId]; ok {
			continue
		} // 如果在白名单则直接跳过

		// 查询是否在slb的后端虚拟服务器组里
		vBkServers, err := GetSlbVserverGroupBackendServer(vGroup.VServerGroupId)
		if err != nil {
			msg.Err = err
			return
		}

		for _, server := range vBkServers {
			if server.ServerId == serverId {

				newSet := make([]backendServer, len(vBkServers))
				// 设置 vServerGroup 里的服务器权重

				// 构建新的权重参数
				for i, _bkServer := range vBkServers {
					newSet[i].ServerId = _bkServer.ServerId
					newSet[i].ServerIp = _bkServer.ServerIp
					newSet[i].Type = _bkServer.Type
					newSet[i].Port = strconv.Itoa(*_bkServer.Port)
					newSet[i].Description = _bkServer.Description
					newSet[i].Weight = strconv.Itoa(_bkServer.Weight)
					if _bkServer.ServerId == serverId {
						newSet[i].Weight = strconv.Itoa(weight)
					}
				}

				log.Debugf("设置slb %s 后端虚拟服务器组: %#v", slb.LoadBalancerName, newSet)
				if err = SetSlbVserverGroup(vGroup.VServerGroupId, newSet); err != nil {
					continue
				}

				msg.Result = append(msg.Result, EcsSetResult{
					ServerId:  serverId,
					SlbId:     slb.LoadBalancerId,
					SlbName:   slb.LoadBalancerName,
					GroupType: VgroupBackServer,
					GroupName: vGroup.VServerGroupName,
					From:      server.Weight,
					To:        weight,
				})

				break
			}
		}
	}

}

// 设置机器的所有负载均衡的权重
// 理论上同一类机器，在它所处的业务线上的权重比例应该是一致的
func SetEcsWeight(serverId string, weight int) (result []EcsSetResult, err error) {
	if weight < 0 {
		weight = 0
	}
	if weight > 100 {
		weight = 100
	}

	slbs, err := GetEcsSlb(serverId)
	if err != nil {
		return result, err
	}

	// 逐个执行slb太慢，用并发提速
	// 这里并发只能到单个 slb级别，如果对一个slb里面的不同组并发设置，阿里云接口会抛 ErrorCode: BackendServer.configuring 这个错误 Message: A previous configuration of the load balancer is pending; please try again later.
	resultChan := make(chan setWeightMsg, len(slbs))
	for _, slb := range slbs {
		go updateEcsSlbWeight(serverId, slb, weight, &resultChan)
	}

	for range slbs {
		msg := <-resultChan
		if msg.Err != nil {
			err = msg.Err
		}
		result = append(result, msg.Result...)
	}

	return
}
