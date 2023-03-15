package utils


import jsoniter "github.com/json-iterator/go"

type ApiJson struct {
	Status int          `json:"status"`
	Msg    string       `json:"msg,omitempty"`
	Data   *interface{} `json:"data,omitempty"`
}

func ApiResource(status int, objects interface{}, msg string) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Msg: msg}
	if objects != nil{
		apijson.Data = &objects
	}
	return
}

var Json = jsoniter.ConfigCompatibleWithStandardLibrary