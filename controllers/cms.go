package controllers

import (
	"github.com/buzhiyun/aliyun-api/cms"
	"github.com/buzhiyun/aliyun-api/ecs"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/buzhiyun/go-utils/log"
	"github.com/kataras/iris/v12"
)

type getDataPointReq struct {
	InstanceId []string `json:"instanceId" example:"i-bp1d1oh9a06r70buf03l"` // 主机Id
	HostName   []string `json:"hostName" example:"sdcf_v3_030" `             // 主机名,支持通配符
}

// SearchHost godoc
// @Summary      获取ECS CPU信息
// @Description  根据ECS  获取CPU使用率
// @Tags         监控CMS
// @Accept       json
// @Produce      json
// @Param   json  body     getDataPointReq   true  "hostName 是 []string"
// @Success      200  {object}   utils.ApiJson{data=[]cms.Datapoint}
// @Failure      400  {object}  utils.ApiJson
// @Failure      500  {object}  utils.ApiJson
// @Router       /api/cms/ecs/cpu [post]
func GetEcsCpu(ctx iris.Context) {
	var data getDataPointReq
	err := ctx.ReadJSON(&data)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	log.Debugf("[cms] %s 尝试刷新查询监控数据", ctx.GetHeader("realip"))

	instanceIds := data.InstanceId

	for _, hostname := range data.HostName {
		for _, instance := range ecs.SearchByName(hostname) {
			instanceIds = append(instanceIds, instance.InstanceId)
		}
	}

	if len(instanceIds) == 0 {
		log.Warnf("无效实例Id或主机名 , %v", data)
		ctx.JSON(utils.ApiResource(200, []cms.Datapoint{}, "ok"))
		return
	}

	resp, err := cms.CMS.GetEcsCpu(instanceIds)
	if err != nil {
		internalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(utils.ApiResource(200, resp, "ok"))
}
