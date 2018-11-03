package controllers

import (
	"encoding/json"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego/logs"
	"strconv"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.StudentID
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	var id models.UserID
	var user models.User
	resp := CommResp{Code: -1}

	err := json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if err != nil {
		logs.Debug("[UserController::Post] invalid json")
		resp.Msg = msgInvalidJSON
		goto Out
	}

	id, err = models.Um.AddUser(&user)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[UserController::Post] AddUser failed")
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = id
Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Param	grade		query 	string	true		"The grade of class"
// @Param	index		query 	string	true		"The number of class"
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	filter := models.Filter{}
	filter.Grade, _ = strconv.Atoi(u.GetString("grade"))
	filter.Index, _ = strconv.Atoi(u.GetString("index"))

	u.Data["json"] = CommResp{
		Code: 0,
		Msg:  msgSuccess,
		Data: models.Um.GetAllUsers(&filter),
	}
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for static block"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	resp := &CommResp{Code: -1}
	tmp := u.GetString(":uid")
	uid, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		logs.Info("invalid user id")
		goto Out
	}

	resp.Data, err = models.Um.GetUser(models.UserID(uid))
	if err != nil {
		resp.Msg = err.Error()
		logs.Debug("getUserID failed")
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	resp := &CommResp{Code: -1}
	tmp := u.GetString(":uid")
	var user models.User
	uid, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		logs.Debug("[] parse uid failed")
		goto Out
	}

	if uid == 0 {
		logs.Debug("invalid uid")
		goto Out
	}

	err = json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	if err != nil {
		resp.Msg = msgInvalidJSON
		goto Out
	}

	err = models.Um.ModUser(&user)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	resp := CommResp{Code: -1}
	tmp := u.GetString(":uid")
	uid, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		logs.Debug("[UserController::Delete] invalid uid")
		goto Out
	}
	err = models.Um.DelUser(models.UserID(uid))
	if err != nil {
		logs.Debug("[UserController::Delete] failed", err)
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *UserController) Login() {
	resp := &CommResp{Code: -1}
	username := u.GetString("username")
	password := u.GetString("password")

	token, err := models.Ac.Login(username, password, models.TypeStudent)
	if err != nil {
		logs.Debug("[UserController::Login] login failed", username, err)
		resp.Msg = err.Error()
		goto Out
	}
	u.SetSession("uid", token)
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = token
Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Param	token		query 	string	true		"The username for login"
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	resp := &CommResp{Code: -1}
	token := u.GetString("username")
	if token == "" {
		logs.Debug("no token")
		resp.Msg = "invalid token"
		goto Out
	}

	if models.Ac.Logout(token) != nil {
		logs.Debug("logout failed")
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess

Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title resetPassword
// @Description reset password of current user
// @Param	token		query 	string	true		"The username for login"
// @Success 200 {string} logout success
// @router /password [post]
func (u *UserController) ResetPassword() {
	resp := &CommResp{Code: -1}
	resp.Code = 0
	resp.Msg = msgSuccess
	u.Data["json"] = resp
	u.ServeJSON()
}
