package controllers

import (
	"encoding/json"

	"github.com/arong/dean/manager"

	"github.com/arong/dean/base"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

// Operations about object
type VoteController struct {
	beego.Controller
}

// @Title GetQuestionnaire
// @Description create object
// @Param	body		body 	voteRequest	true		"The vote info"
// @Success 200 {string} 0
// @Failure 403 body is empty
// @router /survey [post]
func (v *VoteController) GetQuestionnaire() {
	resp := BaseResponse{Code: -1}
	req := manager.GenRequest{}
	ret := models.SurveyPages{}
	var err error

	private := v.Ctx.Input.GetData(base.Private)
	l, ok := private.(manager.LoginInfo)
	if !ok {
		logs.Warn("[VoteController::GetQuestionnaire] bug found")
		resp.Code = base.ErrInternal
		goto Out
	}

	if l.UserType != base.AccountTypeStudent {
		logs.Info("[VoteController::GetQuestionnaire] invalid account type")
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	err = json.Unmarshal(v.Ctx.Input.RequestBody, &req)
	if err != nil {
		logs.Info("[VoteController::GetQuestionnaire] invalid input data", "request", string(v.Ctx.Input.RequestBody))
		resp.Code = base.ErrInvalidInput
		goto Out
	}

	if req.QuestionnaireID == 0 {
		logs.Debug("[VoteController::GetQuestionnaire] invalid questionnaire id")
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	req.StudentID = l.ID

	ret, err = manager.QuestionnaireManager.Generate(req)
	if err != nil {
		logs.Info("[VoteController::GetQuestionnaire] Generate failed", "err", err)
		resp.Code = base.ErrInternal
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = ret
Out:
	v.Data["json"] = resp.Fill()
	v.ServeJSON()
}

// @Title Get
// @Description find object by voteCode
// @Param	voteCode		path	string	true		"the voteCode to verify access"
// @Success 200 {object} models.ScoreInfo
// @router /submit [post]
func (v *VoteController) Submit() {
	resp := BaseResponse{Code: -1}
	var data *models.ClassResp

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = data
	//Out:
	v.Data["json"] = resp
	v.ServeJSON()
}
