package controllers

import (
	"encoding/json"

	"github.com/arong/dean/manager"

	"github.com/arong/dean/base"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Operations about Users
type StudentScoreController struct {
	beego.Controller
}

// @Title Add
// @Description add new record
// @Success 200 {object} base.BaseResponse
// @router /add [post]
func (u *StudentScoreController) Add() {
	request := models.StudentScore{}
	resp := base.BaseResponse{}

	err := json.Unmarshal(u.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Debug("[StudentScoreController::Add] invalid input", "err", err)
		resp.Code = base.ErrInvalidInput
		goto Out
	}

	err = request.Check()
	if err != nil {
		logs.Debug("[StudentScoreController::Add] invalid parameter", "err", err)
		resp.Msg = err.Error()
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	err = manager.SSM.AddRecord(request)
	if err != nil {
		logs.Debug("[StudentScoreController::Add] AddRecord failed", "err", err)
		resp.Code = base.ErrInternal
		resp.Msg = err.Error()
		goto Out
	}

	resp.Msg = msgSuccess
Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router /list [post]
func (u *StudentScoreController) Filter() {
	request := models.StudentFilter{}
	resp := base.BaseResponse{}
	ret := base.CommList{}

	err := json.Unmarshal(u.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Debug("[StudentController::Filter] invalid input", "err", err)
		resp.Code = base.ErrInvalidInput
		goto Out
	}

	err = request.Check()
	if err != nil {
		logs.Debug("[StudentController::Filter] invalid parameter", "err", err)
		resp.Msg = err.Error()
		resp.Code = base.ErrInvalidParameter
		goto Out
	}

	ret = manager.Um.GetAllUsers(request)
	resp.Msg = msgSuccess
	resp.Data = ret
Out:
	u.Data["json"] = resp
	u.ServeJSON()
}
