package controllers

import (
	"encoding/json"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"strconv"
)

// Operations about object
type TeacherController struct {
	beego.Controller
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Teacher	true		"The object content"
// @Success 200 {string} models.Teacher.ID
// @router / [post]
func (o *TeacherController) Post() {
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
// @Description find object by objectid
// @Param	teacherID		path 	string	true		"the teacherID you want to get"
// @Success 200 {object}	models.Teacher
// @Failure 403 :teacherID is empty
// @router /:teacherID [get]
func (o *TeacherController) Get() {
	resp := CommResp{Code: -1}
	var err error
	var id int
	ret := &models.Teacher{}

	teacherID := o.Ctx.Input.Param(":teacherID")
	if teacherID == "" {
		resp.Msg = invalidParam
		goto Out
	}

	id, err = strconv.Atoi(teacherID)
	if err != nil {
		resp.Msg = invalidParam
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

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Teacher
// @router / [get]
func (o *TeacherController) GetAll() {
	resp := &CommResp{
		Code: 0,
		Msg:  msgSuccess,
		Data: models.Tm.GetAll(),
	}
	o.Data["json"] = resp
	o.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (tc*TeacherController) Delete() {
	resp := &CommResp{Code:-1}
	uid := tc.GetString(":uid")
	id, err := strconv.Atoi(uid)
	err = models.Tm.DelTeacher(id)
	if err != nil {

	}
	tc.Data["json"] = resp
	tc.ServeJSON()
}
