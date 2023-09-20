package slb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/buzhiyun/aliyun-api/msg"
	"github.com/buzhiyun/go-utils/log"
	"strings"
	"time"
)

type aclEntry struct {
	Entry   string `json:"entry"`
	Comment string `json:"comment"`
}

func AddIpToAcl(AclId string, IP []string, comment ...string) (err error) {

	entrys := []aclEntry{}
	_c := ""
	if len(comment) > 0 {
		_c = comment[0]
	}
	for _, _ip := range IP {
		if strings.Index(_ip, "/") < 0 {
			_ip = _ip + "/32"
		}
		entrys = append(entrys, aclEntry{
			Entry:   _ip,
			Comment: _c,
		})
	}

	entrysJson, _ := json.MarshalToString(entrys)

	request := slb.CreateAddAccessControlListEntryRequest()
	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(60 * time.Second)

	request.Scheme = "https"

	request.AclId = AclId
	request.AclEntrys = entrysJson

	response, err := client().AddAccessControlListEntry(request)
	if err != nil {
		log.Errorf("添加IP %v 到ACL %s 失败, %s", IP, AclId, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return err
	}

	log.Infof("添加IP %v 到ACL %s 成功 \n%s", IP, AclId, response.GetHttpContentString())
	return
}

func RemoveIpFromAcl(AclId string, IP []string, comment ...string) (err error) {

	entrys := []aclEntry{}
	//_c := ""
	//if len(comment) > 0 {
	//	_c = comment[0]
	//}
	for _, _ip := range IP {
		if strings.Index(_ip, "/") < 0 {
			_ip = _ip + "/32"
		}

		entrys = append(entrys, aclEntry{
			Entry:   _ip,
			Comment: "",
		})
	}

	entrysJson, _ := json.MarshalToString(entrys)

	request := slb.CreateRemoveAccessControlListEntryRequest()
	// 连接超时设置，仅对当前请求有效。
	request.SetConnectTimeout(5 * time.Second)
	// 读超时设置，仅对当前请求有效。
	request.SetReadTimeout(60 * time.Second)

	request.Scheme = "https"

	request.AclId = AclId
	request.AclEntrys = entrysJson

	response, err := client().RemoveAccessControlListEntry(request)
	if err != nil {
		log.Errorf("从ACL %s 删除IP %v 失败, %s", AclId, IP, err.Error())
		msg.AliyunSdkAlert(err.Error())
		return err
	}

	log.Infof("从ACL %s 删除IP %v 成功 \n%s", AclId, IP, response.GetHttpContentString())
	return
}
