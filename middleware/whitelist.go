package middleware

import (
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/buzhiyun/go-utils/cfg"
	"github.com/buzhiyun/go-utils/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func WhiteList(ctx *gin.Context) {
	wl, ok := cfg.Config().GetStrings("security.whitelist")
	if !ok {
		log.Errorf("读取配置 security.whitelist 异常")
		// 这里先放行，免得配置异常造成不可访问
		ctx.Next()
		return
	}

	forwardedIP := ctx.GetHeader("X-Forwarded-For")
	log.Debugf("X-Forwarded-For: [%s]", forwardedIP)
	clientHost := strings.Split(ctx.Request.RemoteAddr, ":")[0]
	// 本机直接放行
	if clientHost == "" || clientHost == "127.0.0.1" {
		ctx.Next()
		return
	}

	var safeIp = false
	for _, ip := range wl {

		// 对走负载均衡 X-Forwarded-For IP 和 直连IP 有一个通过即可
		if utils.MatchWildcard(clientHost, ip) {
			safeIp = true
			ctx.Request.Header.Set("realip", clientHost)
			break
		}

		if utils.MatchWildcard(forwardedIP, ip) {
			safeIp = true
			ctx.Request.Header.Set("realip", forwardedIP)
			break
		}

	}

	if !safeIp {
		log.Warnf("非授权IP %s %s 试图访问 %s", forwardedIP, clientHost, ctx.Request.URL.Path)
		ctx.JSON(http.StatusForbidden, utils.ApiResource(403, nil, ""))
		ctx.Abort()
		return
	}

	ctx.Next()
}
