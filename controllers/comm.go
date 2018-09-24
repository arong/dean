package controllers

import "encoding/json"

type CommResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	defaultResp = `{"code":-1,"msg":"internal error","data":null}`
)

func (cm *CommResp) String() string {
	buff, err := json.Marshal(cm)
	if err != nil {
		return defaultResp
	}
	return string(buff)
}

var (
	invalidJSON = "invalid JSON"
	invalidParam = "invalid parameter"
	msgSuccess = "success"
)