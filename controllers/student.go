package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/arong/dean/base"
	"github.com/arong/dean/manager"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego/logs"

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
	var id int64
	var err error
	user := models.StudentInfo{}
	resp := BaseResponse{Code: -1}

	data := u.Ctx.Input.GetData(base.Data)
	if data == nil {
		resp.Msg = "invalid input"
		logs.Warn("[StudentController::Add] empty input found")
		goto Out
	}

	err = json.Unmarshal(data.(json.RawMessage), &user)
	if err != nil {
		logs.Debug("[StudentController::Add] invalid json", "err", err, "data", string(data.(json.RawMessage)))
		resp.Msg = msgInvalidJSON
		goto Out
	}

	err = user.Check()
	if err != nil {
		logs.Debug("[StudentController::Add] invalid json", "err", err)
		resp.Msg = err.Error()
		goto Out
	}

	id, err = manager.Um.AddStudent(user)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[UserController::Post] AddUser failed")
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = struct {
		ID int64 `json:"id"`
	}{ID: id}

Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Param	grade		query 	string	true		"The grade of class"
// @Param	index		query 	string	true		"The number of class"
// @Success 200 {object} models.User
// @router /filter [post]
func (u *StudentController) Filter() {
	var err error
	request := models.StudentFilter{}
	resp := base.BaseResponse{}
	ret := base.CommList{}

	data, ok := u.Ctx.Input.GetData(base.Data).(json.RawMessage)
	if !ok {
		resp.Msg = "invalid input"
		logs.Warn("[StudentController::Add] empty input found", "data", data)
		goto Out
	}

	logs.Debug("[] ", "data", string(data))
	err = json.Unmarshal(data, &request)
	if err != nil {
		logs.Debug("[StudentController::Filter] invalid input", "err", err)
		resp.Code = base.ErrInvalidInput
		goto Out
	}

	err = request.Check()
	if err != nil {
		logs.Debug("[StudentController::Filter] invalid parameter", "err", err)
		resp.Msg = err.Error()
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	ret = manager.Um.Filter(request)
	resp.Msg = msgSuccess
	resp.Data = ret
Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for static block"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /info/:uid [get]
func (u *StudentController) GetInfo() {
	resp := &BaseResponse{Code: -1}
	tmp := u.GetString(":uid")
	uid, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		logs.Info("invalid user id")
		goto Out
	}

	resp.Data, err = manager.Um.GetStudent(uid)
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
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /update [post]
func (u *StudentController) Update() {
	var err error
	resp := &BaseResponse{Code: -1}
	var user models.StudentInfo

	data, ok := u.Ctx.Input.GetData(base.Data).(json.RawMessage)
	if !ok {
		resp.Msg = "invalid input"
		logs.Warn("[StudentController::Update] empty input found", "data", data)
		goto Out
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		resp.Msg = msgInvalidJSON
		goto Out
	}

	err = manager.Um.UpdateStudent(user)
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

type delStuReq struct {
	IDList models.Int64List `json:"id_list"`
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /delete [post]
func (u *StudentController) Delete() {
	var err error
	resp := BaseResponse{Code: -1}
	request := delStuReq{}
	failed := []int64{}

	data, ok := u.Ctx.Input.GetData(base.Data).(json.RawMessage)
	if !ok {
		resp.Msg = "invalid input"
		logs.Warn("[StudentController::Delete] empty input found", "data", data)
		goto Out
	}

	err = json.Unmarshal(data, &request)
	if err != nil {
		logs.Debug("[StudentController::Delete] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	request.IDList = request.IDList.RemoveZeroNegative()
	if len(request.IDList) == 0 {
		logs.Debug("[StudentController::Delete] empty id list")
		resp.Msg = msgInvalidParam
		goto Out
	}

	failed, err = manager.Um.DelStudent(request.IDList)
	if err != nil {
		logs.Debug("[StudentController::Delete] failed", err)
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	if len(failed) > 0 {
		resp.Data = failed
	}

Out:
	u.Data["json"] = resp
	u.ServeJSON()
}
