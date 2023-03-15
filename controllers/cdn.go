package controllers

import (
	"github.com/buzhiyun/aliyun-api/cdn"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"strings"
)


type RefreshCdnReq struct {
	Urls *[]string `json:"urls,omitempty"  validate:"required" err_info:"必须输入urls []string"` // 主机名,支持通配符
}


func RefreshCdnUrl(ctx iris.Context)  {
	var data RefreshCdnReq
	err := ctx.ReadJSON(&data)
	if err != nil {
		badRequest(ctx,err.Error())
		return
	}

	golog.Infof("[cdn] %s 尝试刷新cdn资源：%s",ctx.GetHeader("realip"),strings.Join(*data.Urls,"\n"))
	resp , err := cdn.RefreshUrl(*data.Urls)
	if err != nil {
		internalServerError(ctx,err.Error())
		return
	}

	ctx.JSON(utils.ApiResource(200,resp, "刷新cdn成功" ))
}
