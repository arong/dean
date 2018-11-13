package controllers

import (
	"encoding/json"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego/logs"
	"strconv"

	"github.com/astaxie/beego"
)

// Operations about Users
type StudentController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.StudentID
// @Failure 403 body is empty
// @router /add [post]
func (u *StudentController) Add() {
	var id models.UserID
	var user models.User
	resp := BaseResponse{Code: -1}

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
func (u *StudentController) GetAll() {
	filter := models.Filter{}
	filter.Grade, _ = strconv.Atoi(u.GetString("grade"))
	filter.Index, _ = strconv.Atoi(u.GetString("index"))

	u.Data["json"] = BaseResponse{
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
// @router /info [get]
func (u *StudentController) GetInfo() {
	resp := &BaseResponse{Code: -1}
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
// @router /update [post]
func (u *StudentController) Update() {
	resp := &BaseResponse{Code: -1}
	tmp := u.GetString(":uid")
	var user models.User
	uid, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		logs.Debug("[StudentController::Update] parse uid failed")
		goto Out
	}

	if uid == 0 {
		logs.Debug("[StudentController::Update] invalid uid")
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
// @router /delete [post]
func (u *StudentController) Delete() {
	resp := BaseResponse{Code: -1}
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
