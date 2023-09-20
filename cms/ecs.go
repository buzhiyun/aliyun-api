package cms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/buzhiyun/go-utils/log"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type reqInstance struct {
	InstanceId string `json:"instanceId"`
}

type Datapoint struct {
	Timestamp  int64   `json:"timestamp"`
	InstanceId string  `json:"instanceId"`
	UserId     string  `json:"userId"`
	Minimum    float64 `json:"Minimum"`
	Maximum    float64 `json:"Maximum"`
	Average    float64 `json:"Average"`
}

func (c *aliyunCms) GetEcsCpu(instanceIds []string) (datapoints []Datapoint, err error) {
	var dimensions []reqInstance

	for _, instanceId := range instanceIds {
		dimensions = append(dimensions, reqInstance{instanceId})
	}

	dimensionsStr, err := json.MarshalToString(dimensions)
	if err != nil {
		log.Errorf("序列化json 数据异常, %s", err.Error())
		return
	}

	request := cms.CreateDescribeMetricLastRequest()

	request.Scheme = "https"

	request.Namespace = "acs_ecs_dashboard"
	request.MetricName = "CPUUtilization"
	//request.Dimensions = "[{\"instanceId\":\"i-bp1d1oh9a06r70buf03h\"}]"
	request.Dimensions = dimensionsStr

	response, err := c.client.DescribeMetricLast(request)
	if err != nil {
		log.Errorf("查询cms datapoints 异常", err.Error())
		return
	}

	err = json.UnmarshalFromString(response.Datapoints, &datapoints)
	if err != nil {
		log.Errorf("反序列化 datapoints 数据异常, %s , %s, %s", err.Error(), response.Message, response.Datapoints)
		return
	}
	return
}
