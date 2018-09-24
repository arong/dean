package controllers

import (
	"strconv"
	"fmt"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"encoding/json"
)

// Operations about object
type TeacherController struct {
	beego.Controller
}


// @Title Create
// @Description create object
// @Param	body		body 	models.Teacher	true		"The object content"
// @Success 200 {string} models.Teacher.ID
// @Failure 403 body is empty
// @router / [post]
func (o *TeacherController) Post() {
	request := voteRequest{}
	resp := commResp{Code: -1}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = invalidJSON
		goto Out
	}

	err = models.Vm.CastVote(request.Scores)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}
Out:
	o.Data["json"] = resp
	o.ServeJSON()
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Teacher
// @Failure 403 :objectId is empty
// @router /:teacherID [get]
func (o *TeacherController) Get() {
	resp := commResp{Code: -1}
	var err error
	var id int
	ret := &models.ScoreInfo{}
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
	fmt.Println("teacherID=", id)
	ret, err = models.Vm.GetScore(id)
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
// @Failure 403 :objectId is empty
// @router / [get]
func (o *TeacherController) GetAll() {
	resp := &commResp{}
	resp.Msg = msgSuccess
	obs := models.Vm.GetAll()
	resp.Data = obs
	o.Data["json"] = resp
	fmt.Println(obs)
	o.ServeJSON()
}
