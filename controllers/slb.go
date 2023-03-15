package controllers

import (
	"github.com/buzhiyun/aliyun-api/ecs"
	"github.com/buzhiyun/aliyun-api/slb"
	"github.com/buzhiyun/aliyun-api/utils"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"strings"
)

type AclListReq struct {
	AclId	string		`json:"acl_id" validate:"required" err_info:"acl_id 不能为空"`
	Host	*[]string	`json:"host,omitempty" validate:"required_without=IP" err_info:"ip 或者 host 不能为空"`
	IP		*[]string	`json:"ip,omitempty" `
	Comment	string		`json:"comment,omitempty"`
}


func AddIpToACL(ctx iris.Context)  {
	var data AclListReq
	err := ctx.ReadJSON(&data)
	if err != nil {
		badRequest(ctx,err.Error())
		return
	}

	var ipList []string

	// 加IP
	if data.IP != nil {
		ipList = append(ipList,*data.IP...)
		return
	}

	// 加 Host 的IP
	if data.Host != nil {
		for _, host := range *data.Host {
			for _, instance := range ecs.SearchByName(host) {
				ipList = append(ipList,instance.IntranetIp)
			}
		}
	}

	golog.Infof("[slb] %s 尝试添加 %s 到 ACL %s",ctx.GetHeader("realip"),strings.Join(ipList,","),data.AclId)
	err = slb.AddIpToAcl(data.AclId,ipList,data.Comment)
	if err != nil {
		internalServerError(ctx,err.Error())
		return
	}

	ctx.JSON(utils.ApiResource(200, nil, "ok"))

}

func DeleteIpFromACL(ctx iris.Context)  {
	var data AclListReq
	err := ctx.ReadJSON(&data)
	if err != nil {
		badRequest(ctx,err.Error())
		return
	}

	var ipList []string

	// 加IP
	if data.IP != nil {
		ipList = append(ipList,*data.IP...)
		return
	}

	// 加 Host 的IP
	if data.Host != nil {
		for _, host := range *data.Host {
			for _, instance := range ecs.SearchByName(host) {
				ipList = append(ipList,instance.IntranetIp)
			}
		}
	}

	golog.Infof("[slb] %s 尝试将 %s 从 ACL %s 中删除",ctx.GetHeader("realip"),strings.Join(ipList,","),data.AclId)
	err = slb.RemoveIpFromAcl(data.AclId,ipList)
	if err != nil {
		internalServerError(ctx,err.Error())
		return
	}

	ctx.JSON(utils.ApiResource(200, nil, "ok"))

}