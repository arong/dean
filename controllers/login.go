package controllers

import (
	"encoding/json"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about Users
type LoginController struct {
	beego.Controller
}

// @Title Login
// @Description Logs user into the system
// @Param	name		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router / [post]
func (l *LoginController) Login() {
	resp := &CommResp{Code: -1}
	req := models.LoginInfo{}
	token := ""

	err := json.Unmarshal([]byte(l.Ctx.Input.RequestBody), &req)
	if err != nil {
		resp.Msg = "invalid data"
		goto Out
	}

	token, err = models.Ac.Login( &req)
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
