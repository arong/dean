package controllers

import "encoding/json"

const (
	defaultResp     = `{"code":-1,"msg":"internal error","data":null}`
	msgInvalidJSON  = "invalid JSON"
	msgInvalidParam = "invalid parameter"
	msgSuccess      = "success"
)

type BaseResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (b *BaseResponse) String() string {
	if b.Msg == "" {
		b.Msg = errMsgMap[b.Code]
	}

	buff, err := json.Marshal(b)
	if err != nil {
		return defaultResp
	}
	return string(buff)
}

func (b *BaseResponse) Fill() *BaseResponse {
	if b.Msg == "" {
		b.Msg = errMsgMap[b.Code]
	}
	return b
}

type CommPage struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type CommList struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

var errMsgMap = map[int]string{
	0: msgSuccess,
}
