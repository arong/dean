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
	var id int64
	var questionnaire models.QuestionnaireInfo
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &questionnaire)
	if err != nil {
		logs.Debug("[QuestionnaireController::Add] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
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

	err = models.Qm.Add(&questionnaire)
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
func (q *QuestionnaireController) Update() {
	var id int64
	var user models.QuestionnaireInfo
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &user)
	if err != nil {
		logs.Debug("[StudentController::Add] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	err = models.Qm.Update(&user)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[UserController::Post] AddUser failed")
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = id
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
	var id int64
	var user models.StudentInfo
	resp := BaseResponse{Code: -1}

	err := json.Unmarshal(q.Ctx.Input.RequestBody, &user)
	if err != nil {
		logs.Debug("[StudentController::Add] invalid json", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	id, err = models.Um.AddUser(&user)
	if err != nil {
		resp.Msg = err.Error()
		logs.Info("[UserController::Post] AddUser failed")
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = id
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
	ret, _ := models.Qm.Filter()

	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = ret
	q.Data["json"] = resp
	q.ServeJSON()
}
