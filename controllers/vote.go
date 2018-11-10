package controllers

import (
	"encoding/json"
	"github.com/arong/dean/models"

	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
)

// Operations about object
type VoteController struct {
	beego.Controller
}

type voteRequest struct {
	VoteCode string
	Scores   []*models.VoteMeta
}

func (vr *voteRequest) Verify() error {
	filter, err := models.Decode(vr.VoteCode)
	if err != nil {
		logs.Warn("[voteRequest::Verify] Decode failed")
		return err
	}

	classResp, err := models.Cm.GetInfo(filter.GetID())
	if err != nil {
		logs.Debug("[voteRequest::Verify] class not found")
		return err
	}

	// check validity
	currMap := make(map[models.UserID]bool)
	for _, v := range classResp.TeacherIDs {
		currMap[v] = true
	}

	// check extra
	for _, v := range vr.Scores {
		if currMap[v.TeacherID] == false {
			return errors.New("access denied")
		}
		delete(currMap, v.TeacherID)
	}

	// check lack
	if len(currMap) > 0 {
		return errors.New("all teacher must be voted")
	}
	return nil
}

// @Title Vote
// @Description create object
// @Param	body		body 	voteRequest	true		"The vote info"
// @Success 200 {string} 0
// @Failure 403 body is empty
// @router / [post]
func (v *VoteController) Post() {
	request := voteRequest{}
	resp := BaseResponse{Code: -1}

	logs.Debug("[VoteController::Post]", "request", string(v.Ctx.Input.RequestBody))
	err := json.Unmarshal(v.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Trace("[VoteController::Post] invalid request format")
		goto Out
	}

	err = request.Verify()
	if err != nil {
		logs.Debug("[VoteController::Post] invalid request", "err", err)
		resp.Msg = err.Error()
		goto Out
	}

	err = models.Vm.CastVote(request.Scores)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = nil
Out:
	v.Data["json"] = resp
	v.ServeJSON()
}

// @Title Get
// @Description find object by voteCode
// @Param	voteCode		path	string	true		"the voteCode to verify access"
// @Success 200 {object} models.ScoreInfo
// @router /:voteCode [get]
func (v *VoteController) Get() {
	resp := BaseResponse{Code: -1}
	var err error
	var filter *models.VoteCodeInfo
	var data *models.ClassResp

	//logs.Debug(v.Ctx)
	voteCode := v.Ctx.Input.Param(":voteCode")
	name := v.GetString("foo")
	logs.Debug("name", name)
	logs.Debug("params", v.Ctx.Input.Params())
	if voteCode == "" {
		resp.Msg = msgInvalidParam
		logs.Debug("no vote code")
		goto Out
	}

	filter, err = models.Decode(voteCode)
	if err != nil {
		resp.Msg = "invalid vote code"
		logs.Debug("invalid vote code")
		goto Out
	}

	logs.Debug("receive a vote", "voteCode", voteCode)

	data, err = models.Cm.GetInfo(filter.GetID())
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
	resp := &BaseResponse{}
	resp.Msg = msgSuccess
	fmt.Println("fuck ")
	obs := models.Vm.GetAll()
	resp.Data = obs
	v.Data["json"] = resp
	fmt.Println(obs)
	v.ServeJSON()
}
