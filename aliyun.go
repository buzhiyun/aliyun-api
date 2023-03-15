package main

import (
	"flag"
	"github.com/buzhiyun/aliyun-api/cdn"
	"github.com/buzhiyun/aliyun-api/controllers"
	"github.com/buzhiyun/aliyun-api/ecs"
	"github.com/buzhiyun/aliyun-api/middleware"
	"github.com/buzhiyun/aliyun-api/slb"
	"github.com/buzhiyun/go-utils/validator"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"strconv"
	"time"
)

type program struct{
	port	int
}

// go-bindata -fs -nomemcopy -prefix "web/dist" ./web/dist/...

func newApp() (app *iris.Application) {
	app = iris.New()
	// OR: basicauth.Default(users)


	app.Validator = validator.New()

	// go get -u github.com/go-bindata/go-bindata/...
	// 静态文件直接打包到程序里  先执行 go-bindata -fs -nomemcopy -prefix "web/dist" ./web/dist/...
	// https://docs.iris-go.com/iris/file-server/http2push-embedded-compression
	//var opts = iris.DirOptions{
	//	IndexName: "index.html",
	//	PushTargetsRegexp: map[string]*regexp.Regexp{
	//		"/": iris.MatchCommonAssets,
	//	},
	//	ShowList: true,
	//	Cache: iris.DirCacheOptions{
	//		Enable:         true,
	//		CompressIgnore: iris.MatchImagesAssets,
	//		Encodings:      []string{"gzip", "deflate", "br", "snappy"},
	//		// Compress files equal or larger than 50 bytes.
	//		CompressMinSize: 50,
	//		Verbose:         1,
	//	},
	//}
	//app.HandleDir("/", AssetFile(), opts)

	//测试的时候允许跨域
	//Cors := cors.New(cors.Options{
	//	AllowedOrigins:   []string{"*"}, // 这里写允许的服务器地址，* 号标识任意
	//	AllowCredentials: true,
	//})


	//api := app.Party("/api", Cors).AllowMethods(iris.MethodOptions)
	api := app.Party("/api")


	// ip白名单
	api.Use(middleware.WhiteList)

	api.PartyFunc("/ecs", func(server iris.Party) {
		server.Post("/search", controllers.SearchHost)
		server.Post("/refresh", controllers.RefreshHost)
		server.Post("/weight", controllers.SetEcsSlbWeight)
	})

	api.PartyFunc("/cdn", func(server iris.Party) {
		server.Post("/refresh",controllers.RefreshCdnUrl)
	})

	slb := api.Party("/slb")

	slb.PartyFunc("/acl" , func(acl iris.Party) {
		acl.Post("/add" , controllers.AddIpToACL)
		acl.Post("/delete" , controllers.DeleteIpFromACL)
	})

	return
}

func (p *program) run() {
	app := newApp()
	app.Run(iris.Addr("0.0.0.0:" + strconv.Itoa(p.port)))
}

func autoRefreshEcs()  {
	for {
		ecs.UpdateEcs()
		time.Sleep(300 * time.Second)
	}
}

func main() {
	golog.SetTimeFormat("2006-01-02 15:04:05.000")

	//if loglevel , ok := cfg.Config().GetString("loglevel") ;ok && loglevel == "debug" {
	//	golog.SetLevel("debug")
	//}
	debug:= flag.Bool("debug",false,"是否开启debug日志")
	port := flag.Int("p",8080,"启动端口")
	flag.Parse()

	if *debug {
		golog.SetLevel("debug")
	}

	if err := ecs.InitECS() ;err !=nil {
		golog.Fatal(err.Error())
	}
	if err := cdn.InitCDN();err !=nil {
		golog.Fatal(err.Error())
	}
	if err := slb.InitSlb() ;err !=nil {
		golog.Fatal(err.Error())
	}

	s := program{*port}

	go autoRefreshEcs()


	s.run()
}