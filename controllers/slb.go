package controllers

import (
	"strings"

	aliyunslb "github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/buzhiyun/aliyun-api/ecs"
	"github.com/buzhiyun/aliyun-api/slb"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/buzhiyun/go-utils/log"
	"github.com/gin-gonic/gin"
)

type AclListReq struct {
	AclId   string    `json:"acl_id" validate:"required" err_info:"acl_id 不能为空"`
	Host    *[]string `json:"host,omitempty" validate:"required_without=IP" err_info:"ip 或者 host 不能为空"`
	IP      *[]string `json:"ip,omitempty" `
	Comment string    `json:"comment,omitempty"`
}

// SearchHost godoc
// @Summary      添加主机到ACL
// @Description  添加主机到ACL
// @Tags         slb
// @Accept       json
// @Produce      json
// @Param   json  body     AclListReq   true  "acl_id 是 slb accessList的ID ； host 和 ip不能同时为空 ，类型均为 []string"
// @Success      200  {object}   utils.ApiJson
// @Failure      400  {object}  utils.ApiJson
// @Failure      500  {object}  utils.ApiJson
// @Router       /api/slb/acl/add [post]
func AddIpToACL(ctx *gin.Context) {
	var data AclListReq
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	var ipList []string

	// 加IP
	if data.IP != nil {
		ipList = append(ipList, *data.IP...)
	}

	// 加 Host 的IP
	if data.Host != nil {
		for _, host := range *data.Host {
			for _, instance := range ecs.SearchByName(host) {
				ipList = append(ipList, ecs.GetInstancesPrivateIP(instance))
			}
		}
	}

	log.Infof("[slb] %s 尝试添加 %s 到 ACL %s", ctx.GetHeader("realip"), strings.Join(ipList, ","), data.AclId)
	err = slb.AddIpToAcl(data.AclId, ipList, data.Comment)
	if err != nil {
		internalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(200, utils.ApiResource(200, nil, "ok"))

}

// SearchHost godoc
// @Summary      从ACL移除主机
// @Description  从ACL移除主机
// @Tags         slb
// @Accept       json
// @Produce      json
// @Param   json  body     AclListReq   true  "acl_id 是 slb accessList的ID ； host 和 ip不能同时为空 ，类型均为 []string"
// @Success      200  {object}   utils.ApiJson
// @Failure      400  {object}  utils.ApiJson
// @Failure      500  {object}  utils.ApiJson
// @Router       /api/slb/acl/delete [post]
func DeleteIpFromACL(ctx *gin.Context) {
	var data AclListReq
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	var ipList []string

	// 加IP
	if data.IP != nil {
		ipList = append(ipList, *data.IP...)
	}

	// 加 Host 的IP
	if data.Host != nil {
		for _, host := range *data.Host {
			for _, instance := range ecs.SearchByName(host) {
				ipList = append(ipList, ecs.GetInstancesPrivateIP(instance))
			}
		}
	}

	log.Infof("[slb] %s 尝试将 %s 从 ACL %s 中删除", ctx.GetHeader("realip"), strings.Join(ipList, ","), data.AclId)
	err = slb.RemoveIpFromAcl(data.AclId, ipList)
	if err != nil {
		internalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(200, utils.ApiResource(200, nil, "ok"))

}

type searchSlbReq struct {
	Slbname string `json:"slbname"  validate:"required_without=Ip" err_info:"slb名称 hostname 和 ip 不能同时为空"` // 主机名,支持通配符
	Ip      string `json:"ip"  validate:""`                                                               // 主机名,支持通配符
	//Fuzzy    *bool  `json:"fuzzy,omitempty" `                                  // 是否模糊搜索 ，默认否
}

// SearchHost godoc
// @Summary      搜索 SLB
// @Description  搜索 SLB
// @Tags         slb
// @Accept       json
// @Produce      json
// @Param   json  body     searchSlbReq   true  "slb名称 hostname 和 ip 不能同时为空"
// @Success      200  {object}   utils.ApiJson
// @Failure      400  {object}  utils.ApiJson
// @Failure      500  {object}  utils.ApiJson
// @Router       /api/slb/search [post]
func SearchSlb(ctx *gin.Context) {
	var data searchSlbReq
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	var slbs []aliyunslb.LoadBalancer

	if data.Slbname != "" {
		slbs, err = slb.SearchByName(data.Slbname)
	}

	if data.Ip != "" {
		slbs, err = slb.SearchByIp(data.Ip)
	}
	if err != nil {
		internalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(200, utils.ApiResource(200, slbs, "ok"))
}

// RefreshSlb godoc
// @Summary      刷新SLB配置
// @Description  刷新SLB配置
// @Tags         slb
// @Accept       json
// @Produce      json
// @Success      200  {object}   utils.ApiJson
// @Failure      400  {object}  utils.ApiJson
// @Failure      500  {object}  utils.ApiJson
// @Router       /api/slb/refresh [post]
func RefreshSlb(ctx *gin.Context) {
	err := slb.RefreshSlb()
	if err != nil {
		internalServerError(ctx, err.Error())
		ctx.JSON(500, utils.ApiResource(500, nil, err.Error()))
		return
	}
	ctx.JSON(200, utils.ApiResource(200, nil, "ok"))
}
