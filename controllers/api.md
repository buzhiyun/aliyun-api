
- ### 接口列表

  目前仅提供给发布系统使用，所以仅有以下接口，参数均由 application/json 传递

  - #### /api/ecs/refresh

    _刷新缓存的ecs列表_
    
    method: **POST**
  
  - #### /api/ecs/search
    
    _搜索主机_
  
    method: POST
  
    | 参数     | 类型   | 是否必须             | 说明                      |
    | -------- | ------ | -------------------- | ------------------------- |
    | hostname | string | 与ip字段二选一       | 主机名称，支持 * 的通配符 |
    | ip       | string | 与hostname字段二选一 | 主机ip地址                |

  - #### /api/ecs/weight
    
    _设置主机在负载均衡上的权重_

    method: **POST**

    | 参数     | 类型   | 是否必须 | 说明                                        |
    | -------- | ------ | -------- | ------------------------------------------- |
    | hostname | string | 是       | 主机名称                                    |
    | weight   | int    | 是       | 该主机在所有负载均衡上的权重值 ，0-100 之间 |

  - #### /api/slb/acl/add
    _负载均衡上的访问控制 添加IP_

    method: **POST**

    | 参数    | 类型     | 是否必须     | 说明                     |
    | ------- | -------- | ------------ | ------------------------ |
    | acl_id  | string   | 是           | slb 上面访问控制的acl id |
    | host    | []string | 与ip二选一   | 主机名称                 |
    | ip      | []string | 与host二选一 | ip 列表                  |
    | comment | string   | 否           | 注释                     |

  - #### /api/slb/acl/delete

    _负载均衡上的访问控制 删除IP_

    method: **POST**

    | 参数    | 类型     | 是否必须     | 说明                     |
    | ------- | -------- | ------------ | ------------------------ |
    | acl_id  | string   | 是           | slb 上面访问控制的acl id |
    | host    | []string | 与ip二选一   | 主机名称                 |
    | ip      | []string | 与host二选一 | ip 列表                  |

  - #### /api/cdn/refresh

    _刷新CDN资源_

    method: **POST**

    | 参数    | 类型     | 是否必须     | 说明        |
    | ------- | -------- | ------------ | ----------- |
    | urls    | []string | 是           | 要刷新的url |
    
    