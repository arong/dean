package controllers

import (
	"encoding/json"
	"github.com/arong/dean/base"
	"strconv"

	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// ClassController manage class
type ClassController struct {
	beego.Controller
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Teacher	true		"The object content"
// @Success 200 {string} models.Teacher.ID
// @router /add [post]
func (c *ClassController) Add() {
	request := models.Class{}
	resp := base.BaseResponse{Code: -1}
	var id int

	logs.Trace("[ClassController::Add]", "request", string(c.Ctx.Input.RequestBody))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[ClassController::Add] invalid json")
		goto Out
	}

	// 班级
	err = request.Check()
	if err != nil {
		resp.Msg = err.Error()
		logs.Debug("[ClassController::Add] invalid parameter")
		goto Out
	}

	if request.MasterID > 0 {
		// check masterID
		if !models.Tm.CheckID(request.MasterID) {
			logs.Debug("[ClassController::Add] invalid head teacher id")
			resp.Msg = "invalid head teacher id"
			goto Out
		}
	}

	// 去重
	request.TeacherList = request.TeacherList.Deduplicate()

	for _, v := range request.TeacherList {
		if v.TeacherID == 0 || v.SubjectID == 0 {
			resp.Code = base.ErrPartialFailed
			resp.Msg = "invalid input"
			goto Out
		}
	}

	if len(request.TeacherList) > 0 {
		err = models.Tm.CheckInstructorList(request.TeacherList)
		if err != nil {
			logs.Debug("[ClassController::Add] invalid instructor list")
			resp.Msg = err.Error()
			goto Out
		}
	}

	id, err = models.Cm.AddClass(&request)
	if err != nil {
		resp.Msg = err.Error()
		logs.Debug("[ClassController::Add] AddClass failed", "err", err)
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = id
Out:
	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @router /update [post]
func (u *ClassController) Update() {
	var class models.Class
	resp := BaseResponse{Code: -1}
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &class)
	if err != nil {
		logs.Debug("[ClassController::Update] invalid json input", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	class.TeacherList = class.TeacherList.Deduplicate()
	err = class.Check()
	if err != nil {
		logs.Debug("[ClassController::Update] invalid data", "err", err)
		resp.Msg = msgInvalidParam
		goto Out
	}

	err = models.Cm.ModifyClass(&class)
	if err != nil {
		logs.Debug("[ClassController::Update] ModifyClass failed", "err", err)
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = nil

Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

type delRequest struct {
	IDList models.ClassIDList `json:"id_list"`
}

// @Title Delete
// @Description delete the user
// @Param	classID		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /delete [post]
func (c *ClassController) Delete() {
	request := delRequest{}
	resp := BaseResponse{Code: -1}
	ret := models.ClassIDList{}

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Debug("[ClassController::Update] invalid json input", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	ret, err = models.Cm.DelClass(request.IDList)
	if err != nil {
		logs.Debug("[ClassController::Delete] failed", "err", err)
		resp.Msg = err.Error()
		goto Out
	}

	if len(ret) > 0 {
		resp.Code = -3
		resp.Msg = "partial failed"
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title Get
// @Description create object
// @Success 200 {object} models.ClassResp
// @router /info [get]
func (c *ClassController) Info() {
	resp := BaseResponse{Code: -1}
	var data *models.Class

	v := c.Ctx.Input.Query("id")
	id, err := strconv.Atoi(v)
	if err != nil {
		logs.Debug("[Info] invalid class id", "id", v)
		resp.Msg = msgInvalidParam
		goto Out
	}

	data, err = models.Cm.GetInfo(id)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = data
Out:
	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Teacher
// @router /filter [get]
func (c *ClassController) Filter() {
	resp := &BaseResponse{
		Code: 0,
		Msg:  msgSuccess,
	}
	tmp := models.Cm.Filter()

	resp.Data = tmp
	c.Data["json"] = resp
	c.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Teacher
// @router /list [get]
func (c *ClassController) GetAll() {
	resp := &BaseResponse{
		Code: 0,
		Msg:  msgSuccess,
	}
	tmp := models.Cm.GetAll()

	resp.Data = tmp
	c.Data["json"] = resp
	c.ServeJSON()
}
