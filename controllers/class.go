package controllers

import (
	"encoding/json"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"strconv"
)

// Operations about object
type ClassController struct {
	beego.Controller
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Teacher	true		"The object content"
// @Success 200 {string} models.Teacher.ID
// @router / [post]
func (o *ClassController) Post() {
	request := models.Teacher{}
	resp := CommResp{Code: -1}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = invalidJSON
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

// @Title Get
// @Description create object
// @Success 200 {object} models.ClassResp
// @router /:grade:index [get]
func (c *ClassController) Get() {
	request := models.Filter{}
	resp := CommResp{Code: -1}
	var data *models.ClassResp
	var err error

	request.Grade,_ = strconv.Atoi(c.GetString(":grade"))
	request.Index,err = strconv.Atoi(c.GetString(":index"))
	if err != nil {
		resp.Msg = invalidJSON
		goto Out
	}

	data, err = models.Cm.GetInfo(&request)
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
		Data: models.Cm.GetAll(),
	}
	c.Data["json"] = resp
	c.ServeJSON()
}
