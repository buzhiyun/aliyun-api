package controllers

import (
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/kataras/iris/v12"
	"net/http"
)

// http 400
func badRequest(ctx iris.Context, errMsg string)  {
	ctx.StatusCode(http.StatusBadRequest)
	ctx.JSON(utils.ApiResource(400, nil ,errMsg))
}


// http 500
func internalServerError(ctx iris.Context, errMsg string)  {
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.JSON(utils.ApiResource(500, nil ,errMsg))
}