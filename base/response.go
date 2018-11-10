package base

import "encoding/json"

type BaseResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type CommPage struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type CommList struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

const (
	defaultResp = `{"code":-1,"msg":"internal error","data":null}`
)

func (cm *BaseResponse) String() string {
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
