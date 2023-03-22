package controllers

import (
	"fmt"
	aliyun_ecs "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/buzhiyun/aliyun-api/ecs"
	"github.com/buzhiyun/aliyun-api/slb"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
)




type searchHostReq struct {
	Hostname string `json:"hostname"  validate:"required_without=Ip" err_info:"主机名称 hostname 和 ip 不能同时为空"` // 主机名,支持通配符
	Ip string `json:"ip"  validate:""` // 主机名,支持通配符
	//Fuzzy    *bool  `json:"fuzzy,omitempty" `                                  // 是否模糊搜索 ，默认否
}



// SearchHost godoc
// @Summary      搜索主机，支持通配符*
// @Description  获取主机列表
// @Tags         ecs
// @Accept       json
// @Produce      json
// @Param   json  body     searchHostReq   true  "hostname 和 ip 不能同时为空"
// @Success      200  {object}   utils.ApiJson
// @Failure      400  {object}  utils.ApiJson
// @Failure      500  {object}  utils.ApiJson
// @Router       /api/ecs/search [post]
func SearchHost(ctx iris.Context)  {
	var data searchHostReq
	err := ctx.ReadJSON(&data)
	if err != nil {
		badRequest(ctx,err.Error())
		return
	}

	var ecsList []aliyun_ecs.Instance
	if len(data.Hostname) > 0 {
		ecsList = ecs.SearchByName(data.Hostname)

		//如果两个条件都传过来了 就取交集，二次筛选
		if len(data.Ip) >0 {
			var _ecsList []aliyun_ecs.Instance
			for _, instance := range ecsList {
				for _, ip := range ecs.GetInstanceIp(instance) {
					if utils.MatchWildcard(ip,data.Ip){
						_ecsList = append(_ecsList, instance)
						break
					}
				}
			}
			ecsList = _ecsList

		}
	} else {
		ecsList = ecs.SearchByIP(data.Ip)
	}

	ctx.JSON(utils.ApiResource(200, ecsList ,"ok"))

}



// SearchHost godoc
// @Summary      刷新ecs实例列表
// @Description  刷新ecs实例列表
// @Tags         ecs
// @Accept       json
// @Produce      json
// @Success      200  {object}   utils.ApiJson
// @Failure      400  {object}  utils.ApiJson
// @Failure      500  {object}  utils.ApiJson
// @Router       /api/ecs/refresh [post]
//
func RefreshHost(ctx iris.Context)  {
	golog.Infof("[ecs] %s 尝试刷新实例列表",ctx.GetHeader("realip"))

	refreshCount , err := ecs.UpdateEcs()
	if err != nil {
		internalServerError(ctx,err.Error())
		return
	}

	ctx.JSON(utils.ApiResource(200, refreshCount,fmt.Sprintf("成功刷新 %v 个实例", refreshCount)))

}


type setEcsWeight struct {
	Hostname string `json:"hostname"  validate:"required" err_info:"主机名称hostname不能为空"` // 主机名,支持通配符
	Weight   *int    `json:"weight"  validate:"required,gte=0,lte=100" err_info:"权重weight 必须是0-100之间的非空值"`       // 权重  *int 防止0在json的时候被丢掉
}


// 会设置服务器 所有负载均衡里的权重
// SearchHost godoc
// @Summary      设置所有包含该服务器的负载均衡里的权重
// @Description  会设置服务器 所有负载均衡里的权重
// @Tags         ecs
// @Accept       json
// @Produce      json
// @Param   json  body     setEcsWeight   true  "hostname 必填 ；weight 为 0-100之间的数字 必填 "
// @Success      200  {object}   utils.ApiJson
// @Failure      400  {object}  utils.ApiJson
// @Failure      500  {object}  utils.ApiJson
// @Router       /api/ecs/search [post]
func SetEcsSlbWeight(ctx iris.Context)  {
	var data setEcsWeight
	err := ctx.ReadJSON(&data)
	if err != nil {
		badRequest(ctx,err.Error())
		return
	}

	golog.Infof("[ecs] %s 尝试设置服务器在负载均衡中的权重",ctx.GetHeader("realip"))

	servers := ecs.SearchByName(data.Hostname)
	var result []slb.EcsSetResult
	var msg string
	for _, server := range servers {
		setResult ,err := slb.SetEcsWeight(server.InstanceId,*data.Weight)
		if err != nil {
			msg = err.Error()
			continue
		}
		result = append(result, setResult...)
	}

	ctx.JSON(utils.ApiResource(200, result, msg))

}