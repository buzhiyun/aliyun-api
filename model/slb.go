package model

import (
	aliyunslb "github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
)

type SlbWithListener struct {
	aliyunslb.LoadBalancer
	Listeners []SlbListener `json:"Listeners,omitempty"`
}

type SlbListener struct {
	AclType           string          `json:"AclType" xml:"AclType"`
	Status            string          `json:"Status" xml:"Status"`
	VServerGroupId    string          `json:"VServerGroupId" xml:"VServerGroupId"`
	ListenerProtocol  string          `json:"ListenerProtocol" xml:"ListenerProtocol"`
	ListenerPort      int             `json:"ListenerPort" xml:"ListenerPort"`
	AclId             string          `json:"AclId" xml:"AclId"`
	Scheduler         string          `json:"Scheduler" xml:"Scheduler"`
	Description       string          `json:"Description" xml:"Description"`
	AclStatus         string          `json:"AclStatus" xml:"AclStatus"`
	BackendServerPort int             `json:"BackendServerPort" xml:"BackendServerPort"`
	BackendProtocol   string          `json:"BackendProtocol" xml:"BackendProtocol"`
	AclIds            []string        `json:"AclIds" xml:"AclIds"`
	Tags              []aliyunslb.Tag `json:"Tags" xml:"Tags"`
}
