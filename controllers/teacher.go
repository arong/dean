package controllers

import (
	"encoding/json"
	"strconv"

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
func (o *TeacherController) Add() {
	request := models.Teacher{}
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[TeacherController] Unmarshal failed", "err", err)
		goto Out
	}

	err = models.Tm.AddTeacher(&request)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = nil
Out:
	o.Data["json"] = resp
	o.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	body		body 	models.Teacher	true		"teacher info"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /modify [post]
func (u *TeacherController) Modify() {
	resp := &BaseResponse{Code: -1}
	var request models.Teacher
	var err error

	err = json.Unmarshal(u.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Info("[TeacherController::Put] unmarshal failed", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	if request.TeacherID == 0 {
		logs.Info("[TeacherController::Put] invalid teacher id")
		goto Out
	}

	err = models.Tm.ModTeacher(&request)
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
	ret := &models.Teacher{}

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

	ret, err = models.Tm.GetTeacherInfo(id)
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
func (o *TeacherController) Filter() {
	resp := &BaseResponse{}
	request := models.TeacherFilter{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &request)
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

	resp.Data = models.Tm.Filter(&request)

Out:
	o.Data["json"] = resp
	o.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.TeacherListResp
// @router /list [get]
func (o *TeacherController) GetAll() {
	o.Data["json"] = &BaseResponse{Msg: msgSuccess, Data: models.Tm.GetAll()}
	o.ServeJSON()
}

type DeleteTeacherReq struct {
	IDList []models.UserID `json:"id_list"`
}

type DeleteTeacherResp struct {
	FailedList []models.UserID
}

// @Title Delete
// @Description delete the user
// @Param	body		body 	models.TeacherFilter	true		"The object content"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /delete [post]
func (tc *TeacherController) Delete() {
	request := DeleteTeacherReq{}
	resp := &BaseResponse{Code: -1}
	ret := DeleteTeacherResp{}

	err := json.Unmarshal(tc.Ctx.Input.RequestBody, &request)
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

	ret.FailedList, err = models.Tm.DelTeacher(request.IDList)
	if err != nil {
		logs.Debug("[TeacherController::Delete] failed", "err", err)
		resp.Msg = err.Error()
		resp.Data = ret
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	tc.Data["json"] = resp
	tc.ServeJSON()
}
