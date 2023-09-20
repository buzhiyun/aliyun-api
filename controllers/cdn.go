package controllers

import (
	"github.com/buzhiyun/aliyun-api/cdn"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/buzhiyun/go-utils/log"
	"github.com/kataras/iris/v12"
	"strings"
)

type RefreshCdnReq struct {
	Urls *[]string `json:"urls,omitempty"  validate:"required" err_info:"必须输入urls []string"` // 主机名,支持通配符
}

// SearchHost godoc
// @Summary      刷新CDN
// @Description  根据提供的URL去刷新CDN
// @Tags         cdn
// @Accept       json
// @Produce      json
// @Param   json  body     RefreshCdnReq   true  "urls 是 []string"
// @Success      200  {object}   utils.ApiJson
// @Failure      400  {object}  utils.ApiJson
// @Failure      500  {object}  utils.ApiJson
// @Router       /api/cdn/refresh [post]
func RefreshCdnUrl(ctx iris.Context) {
	var data RefreshCdnReq
	err := ctx.ReadJSON(&data)
	if err != nil {
		badRequest(ctx, err.Error())
		return
	}

	log.Infof("[cdn] %s 尝试刷新cdn资源：%s", ctx.GetHeader("realip"), strings.Join(*data.Urls, "\n"))
	resp, err := cdn.RefreshUrl(*data.Urls)
	if err != nil {
		internalServerError(ctx, err.Error())
		return
	}

	ctx.JSON(utils.ApiResource(200, resp, "刷新cdn成功"))
}
