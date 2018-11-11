package controllers

import (
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
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [post]
func (l *AuthController) Login() {
	resp := &BaseResponse{Code: -1}
	req := models.LoginInfo{}
	token := ""

	err := json.Unmarshal([]byte(l.Ctx.Input.RequestBody), &req)
	if err != nil {
		resp.Msg = "[Login] invalid request"
		goto Out
	}

	token, err = models.Ac.Login(&req)
	if err != nil {
		logs.Debug("[UserController::Login] login failed", err)
		resp.Msg = err.Error()
		goto Out
	}
	l.SetSession("token", token)
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = token
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
