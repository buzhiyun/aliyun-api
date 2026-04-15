package main

import (
	"flag"
	"strconv"

	"github.com/buzhiyun/aliyun-api/cdn"
	"github.com/buzhiyun/aliyun-api/controllers"
	_ "github.com/buzhiyun/aliyun-api/docs"
	"github.com/buzhiyun/aliyun-api/ecs"
	"github.com/buzhiyun/aliyun-api/middleware"
	"github.com/buzhiyun/aliyun-api/slb"
	"github.com/buzhiyun/go-utils/cfg"
	"github.com/buzhiyun/go-utils/log"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type program struct {
	port int
}

// swag i -g aliyun.go

func newApp() *gin.Engine {
	app := gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	api := app.Group("/api")

	// ip白名单
	api.Use(middleware.WhiteList)

	ecsGroup := api.Group("/ecs")
	{
		ecsGroup.POST("/search", controllers.SearchHost)
		ecsGroup.POST("/refresh", controllers.RefreshHost)
		ecsGroup.POST("/weight", controllers.SetEcsSlbWeight)
	}

	cdnGroup := api.Group("/cdn")
	{
		cdnGroup.POST("/refresh", controllers.RefreshCdnUrl)
	}

	slbGroup := api.Group("/slb")
	slbGroup.POST("/refresh", controllers.RefreshSlb)
	slbGroup.POST("/search", controllers.SearchSlb)
	aclGroup := slbGroup.Group("/acl")
	{
		aclGroup.POST("/add", controllers.AddIpToACL)
		aclGroup.POST("/delete", controllers.DeleteIpFromACL)
	}

	cmsGroup := api.Group("/cms")
	cmsEcsGroup := cmsGroup.Group("/ecs")
	{
		cmsEcsGroup.POST("/cpu", controllers.GetEcsCpu)
		cmsEcsGroup.POST("/mem", controllers.GetEcsMem)
		cmsEcsGroup.POST("/gpu_gpu", controllers.GetEcsGpuGpu)
		cmsEcsGroup.POST("/gpu_mem", controllers.GetEcsGpuMem)
	}

	// swagger 配置
	// 记得运行 swag init -g aliyun.go
	swaggerGroup := app.Group("/swagger")
	swaggerGroup.Use(middleware.WhiteList)

	// 访问 /swagger/ 重定向到 /swagger/index.html，其余走 gin-swagger
	swaggerGroup.GET("/*any", func(c *gin.Context) {
		if c.Param("any") == "/" || c.Param("any") == "" {
			c.Redirect(301, "/swagger/index.html")
			return
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	return app
}

func (p *program) run() {
	app := newApp()
	app.Run("0.0.0.0:" + strconv.Itoa(p.port))
}

func main() {
	if loglevel, ok := cfg.Config().GetString("loglevel"); ok && loglevel == "debug" {
		log.Info("设置日志级别为debug")
		log.SetLevel("debug")
	}
	debug := flag.Bool("debug", false, "是否开启debug日志")
	port := flag.Int("p", 8080, "启动端口")
	flag.Parse()

	if *debug {
		log.SetLevel("debug")
		log.Info("设置日志级别为debug")
	}

	if err := ecs.InitECS(); err != nil {
		log.Fatal(err.Error())
	}
	if err := cdn.InitCDN(); err != nil {
		log.Fatal(err.Error())
	}
	if err := slb.InitSlb(); err != nil {
		log.Fatal(err.Error())
	}

	s := program{*port}

	s.run()
}
