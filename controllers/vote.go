package controllers

import (
	"encoding/json"
	"github.com/arong/dean/models"

	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

// Operations about object
type VoteController struct {
	beego.Controller
}

type voteRequest struct {
	VoteCode string
	Scores   []*models.VoteMeta
}

type verifyVoteCode struct {
	VoteCode string
}

// @Title Create
// @Description create object
// @Param	body		body 	models.TeacherList	true		"The object content"
// @Success 200 {string} models.TeacherList
// @router / [post]
func (o *VoteController) Describe() {
	request := verifyVoteCode{}
	resp := CommResp{Code: -1}
	filter := &models.VoteCodeInfo{}
	data := models.TeacherList{}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = invalidJSON
		goto Out
	}

	filter, err = models.Decode(request.VoteCode)
	if err != nil {
		resp.Msg = "invalid vote code"
		goto Out
	}

	data, err = models.Cm.GetTeacherList(filter.Grade, filter.Index)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = data
Out:
	o.Data["json"] = resp
	o.ServeJSON()
}

// @Title Create
// @Description create object
// @Param	body		body 	voteRequest	true		"The object content"
// @Success 200 {string} 0
// @Failure 403 body is empty
// @router / [post]
func (o *VoteController) Post() {
	request := voteRequest{}
	resp := CommResp{Code: -1}

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
// @Success 200 {object} models.ScoreInfo
// @Failure 403 :objectId is empty
// @router /:teacherID [get]
func (o *VoteController) Get() {
	resp := CommResp{Code: -1}
	var err error
	var id int64
	ret := &models.ScoreInfo{}
	teacherID := o.Ctx.Input.Param(":teacherID")
	if teacherID == "" {
		resp.Msg = invalidParam
		goto Out
	}
	id, err = strconv.ParseInt(teacherID, 10, 64)
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
// @Success 200 {object} models.ScoreInfo
// @Failure 403 :objectId is empty
// @router / [get]
func (o *VoteController) GetAll() {
	resp := &CommResp{}
	resp.Msg = msgSuccess
	obs := models.Vm.GetAll()
	resp.Data = obs
	o.Data["json"] = resp
	fmt.Println(obs)
	o.ServeJSON()
}
