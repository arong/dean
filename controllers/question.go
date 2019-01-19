package controllers

import (
	"encoding/json"

	"github.com/arong/dean/models"

	"github.com/arong/dean/base"
	"github.com/arong/dean/manager"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about object
type QuestionController struct {
	beego.Controller
}

type QuestionFilter struct {
	QuestionnaireID int `json:"questionnaire_id"`
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	manager.User	true		"body for user content"
// @Success 200 {int} manager.User.StudentID
// @Failure 403 body is empty
// @router /add [post]
func (q *QuestionController) Add() {
	var id int
	var question models.QuestionInfo
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &question)
	if err != nil {
		logs.Debug("[QuestionController::Add] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	if question.QuestionnaireID == 0 {
		logs.Debug("[QuestionController::Add] invalid questionnaire id")
		resp.Code = base.ErrInvalidParameter
		resp.Msg = "invalid questionnaire id"
		goto Out
	}

	err = question.Check()
	if err != nil {
		resp.Code = base.ErrInvalidParameter
		resp.Msg = err.Error()
		goto Out
	}

	id, err = manager.QuestionnaireManager.AddQuestion(&question)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[QuestionController::Add] AddUser failed")
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
// @Param	body		body 	manager.User	true		"body for user content"
// @Success 200 {int} manager.User.StudentID
// @Failure 403 body is empty
// @router /update [post]
func (q *QuestionController) Update() {
	var id int64
	var request models.QuestionInfo
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Debug("[QuestionController::Update] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	if request.QuestionID == 0 {
		logs.Debug("[QuestionController::Update] invalid question id")
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	request.Options = request.Options.FilterEmpty()

	err = manager.QuestionnaireManager.UpdateQuestion(&request)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[QuestionController::Update] UpdateQuestion failed")
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = id
Out:
	q.Data["json"] = resp
	q.ServeJSON()
}

// @Title DeleteQuestion
// @Description create users
// @Param	body		body 	manager.User	true		"body for user content"
// @Success 200 {int} manager.User.StudentID
// @Failure 403 body is empty
// @router /delete [post]
func (q *QuestionController) Delete() {
	var id int64
	var request base.SingleID
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Debug("[QuestionController::Delete] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	if request.ID == 0 {
		logs.Debug("[QuestionController::Delete] invalid question id")
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	err = manager.QuestionnaireManager.DeleteQuestion(request.ID)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[QuestionController::Delete] DeleteDeleteQuestion failed")
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = id
Out:
	q.Data["json"] = resp
	q.ServeJSON()
}

// @Title Filter
// @Description find object which meet filter
// @Success 200 {object} manager.ScoreInfo
// @Failure 403 :objectId is empty
// @router /filter [post]
func (q *QuestionController) Filter() {
	resp := BaseResponse{Code: -1}
	request := QuestionFilter{}
	ret := models.QuestionList{}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Warn("[QuestionController::Filter] invalid input data", "request", q.Ctx.Input.RequestBody)
		resp.Code = base.ErrInvalidInput
		goto Out
	}

	if request.QuestionnaireID == 0 {
		logs.Warn("[QuestionController::Filter] invalid id")
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	ret, err = manager.QuestionnaireManager.GetQuestions(request.QuestionnaireID)
	if err != nil {
		logs.Debug("[QuestionController::Filter] GetQuestions failed", "err", err)
		resp.Code = base.ErrInternal
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = ret
Out:
	q.Data["json"] = resp
	q.ServeJSON()
}

// @Title Get Info
// @Description find object which meet filter
// @Success 200 {object} manager.ScoreInfo
// @Failure 403 :objectId is empty
// @router /info [post]
func (q *QuestionController) Info() {
	resp := BaseResponse{Code: -1}
	request := base.SingleID{}
	ret := &models.QuestionInfo{}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Warn("[QuestionController::Info] invalid input data", "request", q.Ctx.Input.RequestBody)
		resp.Code = base.ErrInvalidInput
		goto Out
	}

	if request.ID == 0 {
		logs.Warn("[QuestionController::Info] invalid id")
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	ret, err = manager.QuestionnaireManager.GetQuestionInfo(request.ID)
	if err != nil {
		logs.Debug("[QuestionController::Info] GetQuestionInfo failed", "err", err)
		resp.Code = base.ErrInternal
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = ret
Out:
	q.Data["json"] = resp
	q.ServeJSON()
}
