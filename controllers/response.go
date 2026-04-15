package controllers

import (
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// http 400
func badRequest(ctx *gin.Context, errMsg string)  {
	ctx.JSON(http.StatusBadRequest, utils.ApiResource(400, nil ,errMsg))
}


// http 500
func internalServerError(ctx *gin.Context, errMsg string)  {
	ctx.JSON(http.StatusInternalServerError, utils.ApiResource(500, nil ,errMsg))
}