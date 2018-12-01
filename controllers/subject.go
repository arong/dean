package controllers

import (
	"encoding/json"
	"github.com/arong/dean/base"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about Users
type SubjectController struct {
	beego.Controller
}

// @Title Add
// @Description add new subject
// @Param	grade		query 	string	true		"The grade of class"
// @Success 200 {object} models.User
// @router /add [post]
func (s *SubjectController) Add() {
	resp := base.BaseResponse{}
	request := models.SubjectInfo{}

	err := json.Unmarshal(s.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[ClassController::Add] invalid json")
		goto Out
	}

	if request.Name == "" || request.Key == "" {
		logs.Debug("[ClassController::Add] invalid name")
		resp.Code = base.ErrInvalidParameter
		resp.Msg = "invalid name"
		goto Out
	}

	request.ID, err = models.Sm.Add(request)
	if err != nil {
		logs.Debug("[ClassController::Add] add failed", "err", err)
		resp.Code = base.ErrInvalidParameter
		resp.Msg = err.Error()
		goto Out
	}

Out:
	s.Data["json"] = resp
	s.ServeJSON()
}

// @Title Update
// @Description add new subject
// @Param	grade		query 	string	true		"The grade of class"
// @Success 200 {object} models.User
// @router /update [post]
func (s *SubjectController) Update() {
	resp := base.BaseResponse{}
	request := models.SubjectInfo{}

	err := json.Unmarshal(s.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[ClassController::Update] invalid json")
		goto Out
	}

	if request.ID == 0 {
		logs.Debug("[ClassController::Update] invalid id")
		resp.Code = base.ErrInvalidParameter
		resp.Msg = "invalid id"
		goto Out
	}

	if request.Key == "" {
		logs.Debug("[ClassController::Update] invalid name")
		resp.Code = base.ErrInvalidParameter
		resp.Msg = "invalid name"
		goto Out
	}

	err = models.Sm.Update(request)
	if err != nil {
		logs.Debug("[ClassController::Update] add failed", "err", err)
		resp.Code = base.ErrInvalidParameter
		resp.Msg = err.Error()
		goto Out
	}

Out:
	s.Data["json"] = resp
	s.ServeJSON()
}

// @Title Delete
// @Description delete subject
// @Param	grade		query 	string	true		"The grade of class"
// @Success 200 {object} models.User
// @router /delete [post]
func (s *SubjectController) Delete() {
	request := base.DelList{}
	failedList := []int{}
	resp := base.BaseResponse{}

	err := json.Unmarshal([]byte(s.Ctx.Input.RequestBody), &request)
	if err != nil {
		logs.Debug("[SubjectController::Delete] invalid request", "err", err)
		resp.Code = base.ErrInvalidInput
		goto Out
	}

	for _, v := range request.IDList {
		if v == 0 {
			logs.Debug("[SubjectController::Delete] invalid class id")
			resp.Code = base.ErrInvalidParameter
			goto Out
		}

		err = models.Sm.Delete(v)
		if err != nil {
			logs.Debug("[SubjectController::Delete] Delete failed", "err", err)
			resp.Code = base.ErrPartialFailed
			failedList = append(failedList, v)
		}
	}

Out:
	s.Data["json"] = resp
	s.ServeJSON()
}

// @Title GetAll
// @Description get all Subject
// @Param	grade		query 	string	true		"The grade of class"
// @Success 200 {object} models.User
// @router /list [get]
func (s *SubjectController) GetAll() {
	s.Data["json"] = BaseResponse{
		Code: 0,
		Msg:  msgSuccess,
		Data: models.Sm.GetAll(),
	}
	s.ServeJSON()
}
