package controllers

import (
	"fmt"

	aliyunEci "github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/buzhiyun/aliyun-api/eci"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/gin-gonic/gin"
)

type searchEciReq struct {
	Name string `json:"name"  validate:"required_without=Ip" err_info:"名称 name 和 ip 不能同时为空"` // 主机名,支持通配符
	Ip   string `json:"ip"  validate:""`                                                     // 主机名,支持通配符
	//Fuzzy    *bool  `json:"fuzzy,omitempty" `                                  // 是否模糊搜索 ，默认否
}

// SearchEci ...
// @Tags ECI
// @Summary 搜索ECI实例
// @Description 搜索ECI实例
// @Accept json
// @Produce json
// @Param req body searchEciReq true "request"
// @Success 200 {object} utils.ApiJson{} "response"
// @Router /api/eci/search [post]
func SearchEci(ctx *gin.Context) {
	var req searchEciReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequest(ctx, err.Error())
		return
	}
	var containerGroups []aliyunEci.DescribeContainerGroupsContainerGroup0
	if req.Name != "" {
		containerGroups = eci.Client.SearchContainer(req.Name, eci.SearchName)
	} else if req.Ip != "" {
		containerGroups = eci.Client.SearchContainer(req.Ip, eci.SearchIP)
	}
	ctx.JSON(200, utils.ApiResource(200, containerGroups, "ok"))
}

// refresh Eci ContainerGroup
// @Tags ECI
// @Summary 刷新ECI容器组
// @Description 刷新ECI容器组
// @Accept json
// @Produce json
// @Success 200 {object} utils.ApiJson{} "response"
// @Router /api/eci/refresh [post]
func RefreshEciContainerGroup(ctx *gin.Context) {

	refreshCount, err := eci.Client.RefreshEciContainer()
	if err != nil {
		internalServerError(ctx, err.Error())
		return
	}
	ctx.JSON(200, utils.ApiResource(200, refreshCount, fmt.Sprintf("成功刷新 %v 个实例", refreshCount)))
}
