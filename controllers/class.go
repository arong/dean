package controllers

import (
	"encoding/json"
	"sort"
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
// @router / [post]
func (c *ClassController) Post() {
	request := models.Class{}
	resp := CommResp{Code: -1}
	var id models.ClassID

	logs.Trace("[ClassController::Post]", "request", string(c.Ctx.Input.RequestBody))
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[ClassController::Post] invalid json")
		goto Out
	}

	id, err = models.Cm.AddClass(&request)
	if err != nil {
		resp.Msg = err.Error()
		logs.Debug("[ClassController::Post] AddClass failed", "err", err)
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
// @router / [put]
func (u *ClassController) Put() {
	var class models.Class
	resp := CommResp{Code: -1}
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &class)
	if err != nil {
		logs.Debug("[ClassController::Put] invalid json input", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	err = models.Cm.ModifyClass(&class)
	if err != nil {
		logs.Debug("[lassController::Put] ModifyClass failed", "err", err)
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

// @Title Delete
// @Description delete the user
// @Param	classID		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:classID [delete]
func (c *ClassController) Delete() {
	resp := CommResp{Code: -1}
	classID, err := strconv.Atoi(c.GetString(":classID"))
	if err != nil {
		logs.Debug("[ClassController::Delete] invalid input param")
		resp.Msg = msgInvalidParam
		goto Out
	}

	err = models.Cm.DelClass(models.ClassID(classID))
	if err != nil {
		logs.Debug("[ClassController::Delete] failed", "err", err)
		resp.Msg = err.Error()
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
// @router /:classID [get]
func (c *ClassController) Get() {
	resp := CommResp{Code: -1}
	var data *models.ClassResp
	var err error

	id, err := strconv.Atoi(c.GetString(":classID"))
	if err != nil {
		resp.Msg = msgInvalidJSON
		goto Out
	}

	data, err = models.Cm.GetInfo(models.ClassID(id))
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
// @router / [get]
func (c *ClassController) GetAll() {
	resp := &CommResp{
		Code: 0,
		Msg:  msgSuccess,
	}
	tmp := models.Cm.GetAll()
	sort.Sort(tmp)

	resp.Data = CommList{RecordCount: len(tmp), RecordList: tmp}
	c.Data["json"] = resp
	c.ServeJSON()
}
