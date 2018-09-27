package controllers

import (
	"encoding/json"
	"github.com/arong/dean/models"

	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about object
type VoteController struct {
	beego.Controller
}

type voteRequest struct {
	VoteCode string
	Scores   []*models.VoteMeta
}

// @Title Create
// @Description create object
// @Param	body		body 	voteRequest	true		"The object content"
// @Success 200 {string} 0
// @Failure 403 body is empty
// @router / [post]
func (v *VoteController) Post() {
	request := voteRequest{}
	resp := CommResp{Code: -1}

	err := json.Unmarshal(v.Ctx.Input.RequestBody, &request)
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
	v.Data["json"] = resp
	v.ServeJSON()
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.ScoreInfo
// @Failure 403 :objectId is empty
// @router /:voteCode [get]
func (v *VoteController) Get() {
	resp := CommResp{Code: -1}
	var err error
	var filter *models.VoteCodeInfo
	var data *models.ClassResp

	voteCode := v.Ctx.Input.Param(":voteCode")
	if voteCode == "" {
		resp.Msg = invalidParam
		goto Out
	}

	filter, err = models.Decode(voteCode)
	if err != nil {
		resp.Msg = "invalid vote code"
		goto Out
	}

	logs.Debug("receive a vote", "voteCode", voteCode)

	data, err = models.Cm.GetInfo(&filter.Filter)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = data
Out:
	v.Data["json"] = resp
	v.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.ScoreInfo
// @Failure 403 :objectId is empty
// @router / [get]
func (v *VoteController) GetAll() {
	resp := &CommResp{}
	resp.Msg = msgSuccess
	obs := models.Vm.GetAll()
	resp.Data = obs
	v.Data["json"] = resp
	fmt.Println(obs)
	v.ServeJSON()
}
