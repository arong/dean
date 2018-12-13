package controllers

import (
	"encoding/json"

	"github.com/arong/dean/base"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about object
type QuestionnaireController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.StudentID
// @Failure 403 body is empty
// @router /add [post]
func (q *QuestionnaireController) Add() {
	var id int
	var questionnaire models.QuestionnaireInfo
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &questionnaire)
	if err != nil {
		logs.Debug("[QuestionnaireController::Add] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	{
		private := q.Ctx.Input.GetData(base.Private)
		if l, ok := private.(models.LoginInfo); ok {
			questionnaire.Editor = l.LoginName
		} else {
			logs.Debug("[QuestionnaireController::Add] invalid user info", "private", private)
		}
	}

	if questionnaire.QuestionnaireID != 0 {
		logs.Debug("[QuestionnaireController::Add] invalid questionnaire id")
		resp.Code = base.ErrInvalidParameter
		resp.Msg = "no id shall be specified"
		goto Out
	}

	err = questionnaire.Check()
	if err != nil {
		resp.Code = base.ErrInvalidParameter
		resp.Msg = err.Error()
		goto Out
	}

	id, err = models.QuestionnaireManager.Add(&questionnaire)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[QuestionnaireController::Add] AddUser failed", "err", err)
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = id
Out:
	q.Data["json"] = resp
	q.ServeJSON()
}

// @Title Update
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.StudentID
// @Failure 403 body is empty
// @router /update [post]
func (q *QuestionnaireController) Update() {
	var request models.QuestionnaireInfo
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Debug("[StudentController::Update] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	if request.QuestionnaireID == 0 {
		logs.Debug("[StudentController::Update] invalid questionnaire id")
		resp.Code = base.ErrInvalidParameter
		resp.Msg = "id shall be specified"
		goto Out
	}

	err = request.Check()
	if err != nil {
		logs.Info("[StudentController::Update] Check failed", "err", err)
		resp.Code = base.ErrInvalidParameter
		resp.Msg = err.Error()
		goto Out
	}

	err = models.QuestionnaireManager.Update(&request)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[StudentController::Update] Update failed", "err", err)
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	q.Data["json"] = resp
	q.ServeJSON()
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.StudentID
// @Failure 403 body is empty
// @router /delete [post]
func (q *QuestionnaireController) Delete() {
	var request base.SingleID
	resp := BaseResponse{}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Debug("[QuestionnaireController::Delete] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		resp.Code = base.ErrInvalidInput
		goto Out
	}

	if request.ID <= 0 {
		resp.Code = base.ErrInvalidParameter
		resp.Msg = "invalid id"
		goto Out
	}

	err = models.QuestionnaireManager.Delete(request.ID)
	if err != nil {
		logs.Info("[QuestionnaireController::Delete] Delete failed", "err", err)
		resp.Code = base.ErrInternal
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	q.Data["json"] = resp
	q.ServeJSON()
}

// @Title Get
// @Description find object which meet filter
// @Success 200 {object} models.ScoreInfo
// @Failure 403 :objectId is empty
// @router /filter [post]
func (q *QuestionnaireController) Filter() {
	resp := BaseResponse{Code: -1}
	ret, _ := models.QuestionnaireManager.Filter()

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = ret
	q.Data["json"] = resp
	q.ServeJSON()
}

// @Title Get
// @Description find object which meet filter
// @Success 200 {object} models.ScoreInfo
// @Failure 403 :objectId is empty
// @router /submit [post]
func (q *QuestionnaireController) Submit() {
	resp := BaseResponse{}
	request := models.QuestionnaireSubmit{}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Code = base.ErrInvalidInput
		resp.Msg = "[QuestionnaireController::Submit] invalid data format"
		goto Out
	}

	{
		p := q.Ctx.Input.GetData(base.Private)
		if l, ok := p.(models.LoginInfo); ok {
			if l.UserType != models.TypeStudent {
				goto Out
			}
			request.StudentID = l.ID
		}
	}

	if request.StudentID == 0 {
		logs.Debug("[QuestionnaireController::Submit] invalid student id")
		resp.Code = base.ErrInvalidParameter
		resp.Msg = "invalid student id"
		goto Out
	}

	err = request.Check()
	if err != nil {
		logs.Debug("[QuestionnaireController::Submit] Check failed", "err", err)
		resp.Code = base.ErrInvalidParameter
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	q.Data["json"] = resp
	q.ServeJSON()
}
