## aliyun-api

简化sdk操作，主要为内部脚本提供接口快速调用阿里云控制台的功能 ,目前阶段主要提供给线上发布脚本使用


- #### 配置文件
  config.yaml：

  文件路径: 与主程序同级或者放在在主程序同级的config目录中。会自动搜索加载
  ```yaml
  aliyun:
    region: "cn-hangzhou"
    key: "xxxxxxxxxxxxxxx"
    secret: "xxxxxxxxxxxxxxxxxxxxxxxxx"
  
  security:
    whitelist:  # 白名单列表 list ，ip支持通配符* ，白名单内的ip可以访问，以后有时间了再改成其他的访问控制
      - "192.168.8.*"

  ```



- #### 运行方式
  - docker build 之后，docker run;
  - go build aliyun.go 之后，直接运行


- #### 参数配置

  | 启动参数 | 默认值 | 说明              |
  | -------- | ------ | ----------------- |
  | -port    | 8080   | 监听的端口号      |
  | -debug   | false  | 是否显示debug日志 |


- #### 接口列表
  详情见 [api.md](controllers/api.md)