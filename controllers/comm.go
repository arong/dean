package controllers

import "encoding/json"

type CommResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type CommList struct {
	RecordCount int         `json:"record_count"`
	RecordList  interface{} `json:"record_list"` //  intended for
}

const (
	defaultResp = `{"code":-1,"msg":"internal error","data":null}`
)

func (cm *CommResp) String() string {
	if cm.Msg == "" {
		cm.Msg = errMsgMap[cm.Code]
	}

	buff, err := json.Marshal(cm)
	if err != nil {
		return defaultResp
	}
	return string(buff)
}

var (
	msgInvalidJSON  = "invalid JSON"
	msgInvalidParam = "invalid parameter"
	msgSuccess      = "success"
)

var errMsgMap = map[int]string{
	0: msgSuccess,
}
