package msg

import (
	"fmt"
	"github.com/buzhiyun/go-utils/cfg"
	"github.com/buzhiyun/go-utils/log"
	"os"
)

func AliyunSdkAlert(errMsg string) {
	toUser, ok := cfg.Config().GetStrings("notice.to_users")
	if !ok || len(toUser) == 0 {
		log.Warnf("获取配置 notice.to_users 异常")
		return
	}
	hostname, _ := os.Hostname()

	sendWechatWorkAppMessage(
		fmt.Sprintf("### 阿里云SDK异常\n> **ERR**: <font color=\"comment\">%s</font>\n**Hostname**: %s", errMsg, hostname),
		toUser)
}
