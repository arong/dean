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

func (cp CommPage) GetRange() (int, int) {
	if cp.Page <= 0 || cp.Size <= 0 {
		return 0, 0
	}
	start := (cp.Page - 1) * cp.Size
	end := cp.Page * cp.Size
	return start, end
}

type CommList struct {
	Total int         `json:"total"`
	List  interface{} `json:"list,omitempty"`
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
