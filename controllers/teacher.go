package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/arong/dean/base"

	"github.com/arong/dean/manager"

	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about object
type TeacherController struct {
	beego.Controller
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Teacher	true		"The object content"
// @Success 200 {string} models.Teacher.TeacherID
// @router /add [post]
func (t *TeacherController) Add() {
	var err error
	var id int64
	request := models.Teacher{}
	resp := BaseResponse{Code: -1}

	data, ok := t.Ctx.Input.GetData(base.Data).(json.RawMessage)
	if !ok {
		resp.Msg = "invalid input"
		logs.Warn("[TeacherController::Add] empty input found", "data", data)
		goto Out
	}

	err = json.Unmarshal(data, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[TeacherController::Add] Unmarshal failed", "err", err)
		goto Out
	}

	id, err = manager.Tm.AddTeacher(&request)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = struct {
		ID int64 `json:"id"`
	}{ID: id}
Out:
	t.Data["json"] = resp
	t.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	body		body 	models.Teacher	true		"teacher info"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /modify [post]
func (t *TeacherController) Modify() {
	resp := &BaseResponse{Code: -1}
	var request models.Teacher
	var err error

	data, ok := t.Ctx.Input.GetData(base.Data).(json.RawMessage)
	if !ok {
		resp.Msg = "invalid input"
		logs.Warn("[TeacherController::Put] empty input found", "data", data)
		goto Out
	}

	err = json.Unmarshal(data, &request)
	if err != nil {
		logs.Info("[TeacherController::Put] unmarshal failed", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	if request.TeacherID == 0 {
		logs.Info("[TeacherController::Put] invalid teacher id")
		goto Out
	}

	err = manager.Tm.UpdateTeacher(&request)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	t.Data["json"] = resp
	t.ServeJSON()
}

// @Title Get
// @Description find object by teacherID
// @Param	teacherID		path 	string	true		"the teacherID you want to get"
// @Success 200 {object}	models.Teacher
// @Failure 403 :teacherID is empty
// @router /info/:teacherID [get]
func (o *TeacherController) Get() {
	resp := BaseResponse{Code: -1}
	var err error
	var id int64
	ret := manager.TeacherInfoResp{}

	teacherID := o.Ctx.Input.Param(":teacherID")
	if teacherID == "" {
		resp.Msg = msgInvalidParam
		goto Out
	}

	id, err = strconv.ParseInt(teacherID, 10, 64)
	if err != nil {
		resp.Msg = msgInvalidParam
		goto Out
	}

	ret, err = manager.Tm.GetTeacherInfo(id)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = ret
Out:
	o.Data["json"] = resp
	o.ServeJSON()
}

// @Title Filter
// @Description get all objects
// @Param	body		body 	models.TeacherFilter	true		"The object content"
// @Success 200 {object} models.TeacherListResp
// @router /filter [post]
func (t *TeacherController) Filter() {
	resp := &BaseResponse{}
	request := models.TeacherFilter{}
	var err error
	data, ok := t.Ctx.Input.GetData(base.Data).(json.RawMessage)
	if !ok {
		resp.Msg = "invalid input"
		logs.Warn("[TeacherController::Put] empty input found", "data", data)
		goto Out
	}

	err = json.Unmarshal(data, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[TeacherController::GetAll] Unmarshal failed", "err", err)
		goto Out
	}

	if request.Page <= 0 || request.Size <= 0 {
		resp.Msg = msgInvalidParam
		logs.Debug("[TeacherController::GetAll] invalid size")
		goto Out
	}

	resp.Data = manager.Tm.Filter(request)

Out:
	t.Data["json"] = resp
	t.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.TeacherListResp
// @router /list [get]
func (o *TeacherController) GetAll() {
	o.Data["json"] = &BaseResponse{Msg: msgSuccess, Data: manager.Tm.GetAll()}
	o.ServeJSON()
}

type DeleteTeacherReq struct {
	IDList []int64 `json:"id_list"`
}

type DeleteTeacherResp struct {
	FailedList []int64
}

// @Title Delete
// @Description delete the user
// @Param	body		body 	models.TeacherFilter	true		"The object content"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /delete [post]
func (t *TeacherController) Delete() {
	request := DeleteTeacherReq{}
	resp := &BaseResponse{Code: -1}
	ret := DeleteTeacherResp{}

	var err error
	data, ok := t.Ctx.Input.GetData(base.Data).(json.RawMessage)
	if !ok {
		resp.Msg = "invalid input"
		logs.Warn("[TeacherController::Put] empty input found", "data", data)
		goto Out
	}

	err = json.Unmarshal(data, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[TeacherController::Delete] Unmarshal failed", "err", err)
		goto Out
	}

	if len(request.IDList) == 0 {
		resp.Msg = "invalid request"
		logs.Debug("[TeacherController::Delete] invalid request")
		goto Out
	}

	ret.FailedList, err = manager.Tm.DelTeacher(request.IDList)
	if err != nil {
		logs.Debug("[TeacherController::Delete] failed", "err", err)
		resp.Msg = err.Error()
		resp.Data = ret
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	if len(ret.FailedList) > 0 {
		resp.Data = struct {
			FailedList []int64 `json:"failed_list,omitempty"`
		}{FailedList: ret.FailedList}
	}
Out:
	t.Data["json"] = resp
	t.ServeJSON()
}
