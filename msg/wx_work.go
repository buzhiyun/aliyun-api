package msg

import (
	"github.com/buzhiyun/go-utils/cfg"
	"github.com/buzhiyun/go-utils/http"
	"github.com/buzhiyun/go-utils/log"
	"strings"
)

/*
发送 MarkDown 消息
https://work.weixin.qq.com/api/doc/90000/90135/90236#markdown%E6%B6%88%E6%81%AF
支持的 Markdown 语法：https://work.weixin.qq.com/api/doc/90000/90135/90236#%E6%94%AF%E6%8C%81%E7%9A%84markdown%E8%AF%AD%E6%B3%95
*/
func sendWechatWorkAppMessage(markdownContent string, toUsers []string) (err error) {

	msgApi, ok := cfg.Config().GetString("notice.msg_api")
	if !ok || msgApi == "" {

		log.Warnf("获取配置 msg_api 异常")
		return
	}

	log.Infof("向 %v 发送企业微信应用消息", toUsers)

	if len(toUsers) == 0 {
		toUsers = []string{"177"}
	}

	var data = struct {
		ToUser  string `json:"touser"`
		Content string `json:"content"`
	}{
		ToUser:  strings.Join(toUsers, "|"),
		Content: markdownContent,
	}

	resp, err := http.HttpPostJson(msgApi+"/api/wechatwork/msg/markdown", data)
	if err != nil {
		log.Errorf("发送企业微信应用消息错误: %s", err.Error())
	} else {
		log.Infof("发送企业微信应用消息返回: %s", resp)
	}
	return
}
