package slb

import (
	"github.com/buzhiyun/go-utils/log"
	jsoniter "github.com/json-iterator/go"
	"strconv"
)

const (
	BackServer = iota
	VgroupBackServer
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type EcsSetResult struct {
	ServerId  string
	SlbId     string
	SlbName   string
	GroupType int // vgroup
	GroupName string
	From      int
	To        int
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

	for _, slb := range slbs {
		// 查询是否在slb的后端服务器里
		bkServers, err := GetSlbBackendServer(slb.LoadBalancerId)
		if err != nil {
			return result, err
		}
		for _, server := range bkServers {
			if server.ServerId == serverId {

				// 后端服务器设置权重
				newSet := make([]backendServer, len(bkServers))
				// 设置 vServerGroup 里的服务器权重

				for i, _bkServer := range bkServers {
					newSet[i].ServerId = _bkServer.ServerId
					newSet[i].ServerIp = _bkServer.ServerIp
					newSet[i].Type = _bkServer.Type
					newSet[i].Description = _bkServer.Description
					newSet[i].Weight = strconv.Itoa(_bkServer.Weight)
					if _bkServer.ServerId == serverId {
						newSet[i].Weight = strconv.Itoa(weight)
					}
				}

				if err = SetSlbBackendServer(slb.LoadBalancerId, newSet); err != nil {
					continue
				}

				result = append(result, EcsSetResult{
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
			return result, err
		}

		for _, vGroup := range vGroups {
			vBkServers, err := GetSlbVserverGroupBackendServer(vGroup.VServerGroupId)
			if err != nil {
				return result, err
			}

			for _, server := range vBkServers {
				if server.ServerId == serverId {

					newSet := make([]backendServer, len(vBkServers))
					// 设置 vServerGroup 里的服务器权重

					for i, _bkServer := range vBkServers {
						newSet[i].ServerId = _bkServer.ServerId
						newSet[i].ServerIp = _bkServer.ServerIp
						newSet[i].Type = _bkServer.Type
						newSet[i].Port = strconv.Itoa(_bkServer.Port)
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

					result = append(result, EcsSetResult{
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

	return
}
