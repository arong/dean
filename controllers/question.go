package controllers

import (
	"encoding/json"

	"github.com/arong/dean/base"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about object
type QuestionController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.StudentID
// @Failure 403 body is empty
// @router /add [post]
func (q *QuestionController) Add() {
	var id int64
	var question models.QuestionInfo
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &question)
	if err != nil {
		logs.Debug("[QuestionnaireController::Add] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	if question.QuestionnaireID == 0 {
		logs.Debug("[QuestionnaireController::Add] invalid questionnaire id")
		resp.Code = base.ErrInvalidParameter
		resp.Msg = "no id shall be specified"
		goto Out
	}

	err = question.Check()
	if err != nil {
		resp.Code = base.ErrInvalidParameter
		resp.Msg = err.Error()
		goto Out
	}

	err = models.QuestionnaireManager.AddQuestion(&question)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[QuestionnaireController::Add] AddUser failed")
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
func (q *QuestionController) Update() {
	var id int64
	var questionInfo models.QuestionInfo
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &questionInfo)
	if err != nil {
		logs.Debug("[QuestionController::Update] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	if questionInfo.QuestionID == 0 {
		logs.Debug("[QuestionController::Update] invalid question id")
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	err = models.QuestionnaireManager.UpdateQuestion(&questionInfo)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[QuestionController::Update] AddUser failed")
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = id
Out:
	q.Data["json"] = resp
	q.ServeJSON()
}

type delReq struct {
	QuestionID int
}

// @Title DeleteQuestion
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.StudentID
// @Failure 403 body is empty
// @router /delete [post]
func (q *QuestionController) Delete() {
	var id int64
	var user delReq
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &user)
	if err != nil {
		logs.Debug("[QuestionController::Delete] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	if user.QuestionID == 0 {
		logs.Debug("[QuestionController::Delete] invalid question id")
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	err = models.QuestionnaireManager.DeleteQuestion(user.QuestionID)
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
// @Success 200 {object} models.ScoreInfo
// @Failure 403 :objectId is empty
// @router /filter [post]
func (q *QuestionController) Filter() {
	resp := BaseResponse{Code: -1}
	request := base.SingleID{}
	ret := models.QuestionList{}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Warn("[QuestionController::Filter] invalid input data", "request", q.Ctx.Input.RequestBody)
		resp.Code = base.ErrInvalidInput
		goto Out
	}

	if request.ID == 0 {
		logs.Warn("[QuestionController::Filter] invalid id")
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	ret, err = models.QuestionnaireManager.GetQuestions(request.ID)
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
