package controllers

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/arong/dean/base"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about Users
type AuthController struct {
	beego.Controller
}

// @Title Login
// @Description Logs user into the system
// @Success 200 {object} models.BaseResponse
// @Failure 403 user not exist
// @router /login [post]
func (l *AuthController) Login() {
	resp := &BaseResponse{Code: -1}
	req := models.LoginRequest{}
	token := ""

	err := json.Unmarshal([]byte(l.Ctx.Input.RequestBody), &req)
	if err != nil {
		logs.Debug("[AuthController::Login] invalid input data")
		resp.Msg = "invalid request"
		goto Out
	}

	err = req.Check()
	if err != nil {
		logs.Info("[UserController::Login] invalid request parameter", "err", err)
		resp.Msg = err.Error()
		goto Out
	}

	token, err = models.Ac.Login(&req)
	if err != nil {
		logs.Debug("[AuthController::Login] login failed", err)
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = struct {
		Token string `json:"token"`
	}{Token: token}
	logs.Info("[AuthController::Login] login success", req.LoginName, resp)

Out:
	l.Data["json"] = resp
	l.ServeJSON()
}

// @Title Update
// @Description update password
// @Success 200 {object} models.BaseResponse
// @Failure 403 user not exist
// @router /update [post]
func (l *AuthController) Update() {
	resp := &BaseResponse{Code: -1}
	req := models.UpdateRequest{}

	err := json.Unmarshal([]byte(l.Ctx.Input.RequestBody), &req)
	if err != nil {
		resp.Msg = "invalid request"
		resp.Code = base.ErrInvalidInput
		goto Out
	}

	// password is sha256 encoded text
	if req.Password == "" || len(req.Password) != sha256.BlockSize {
		logs.Debug("[UserController::Update] login failed", err)
		resp.Msg = "invalid password"
		goto Out
	}

	{
		loginInfo, ok := l.Ctx.Input.GetData(base.Private).(models.LoginInfo)
		if !ok {
			logs.Warn("[UserController::Update] bug found")
			resp.Code = base.ErrInternal
			goto Out
		}
		req.LoginName = loginInfo.LoginName
		req.UserType = loginInfo.UserType
	}

	err = models.Ac.Update(&req)
	if err != nil {
		logs.Debug("[UserController::Update] Update failed", err)
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess

Out:
	l.Data["json"] = resp.Fill()
	l.ServeJSON()
}

// @Title Reset
// @Description reset password
// @Success 200 {object} models.BaseResponse
// @Failure 403 user not exist
// @router /reset [post]
func (l *AuthController) Reset() {
	resp := &BaseResponse{Code: -1}
	req := models.ResetPassReq{}

	err := json.Unmarshal([]byte(l.Ctx.Input.RequestBody), &req)
	if err != nil {
		resp.Msg = "[UserController::Reset] invalid request"
		goto Out
	}

	{
		l, ok := l.Ctx.Input.GetData(base.Private).(models.LoginInfo)
		if !ok {
			resp.Code = base.ErrInternal
			goto Out
		}
		if l.UserType != base.AccountTypeTeacher {
			resp.Code = base.ErrInvalidParameter
			resp.Msg = "permission denied"
			goto Out
		}
	}

	if len(req.Password) != sha256.BlockSize {
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	err = models.Ac.ResetAllStudentPassword(&req)
	if err != nil {
		logs.Debug("[UserController::Reset] login failed", err)
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	logs.Info("[AuthController::Reset] all password reset")

Out:
	l.Data["json"] = resp
	l.ServeJSON()
}

// @Title Logout
// @Description Logs user into the system
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /logout [post]
func (l *AuthController) Logout() {
	resp := &BaseResponse{Code: -1}
	req := base.BaseRequest{}

	err := json.Unmarshal([]byte(l.Ctx.Input.RequestBody), &req)
	if err != nil {
		logs.Debug("[UserController::Login] invalid data")
		resp.Msg = "invalid data"
		goto Out
	}

	err = models.Ac.Logout(req.Token)
	if err != nil {
		logs.Debug("[UserController::Login] login failed", err)
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	l.Data["json"] = resp
	l.ServeJSON()
}
