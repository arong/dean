package controllers

import (
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
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
	s.Data["json"] = BaseResponse{
		Code: 0,
		Msg:  msgSuccess,
		Data: models.Sm.GetAll(),
	}
	s.ServeJSON()
}

// @Title Delete
// @Description delete subject
// @Param	grade		query 	string	true		"The grade of class"
// @Success 200 {object} models.User
// @router /del [post]
func (s *SubjectController) Delete() {
	s.Data["json"] = BaseResponse{
		Code: 0,
		Msg:  msgSuccess,
		Data: models.Sm.GetAll(),
	}
	s.ServeJSON()
}

// @Title GetAll
// @Description get all Subject
// @Param	grade		query 	string	true		"The grade of class"
// @Success 200 {object} models.User
// @router / [get]
func (s *SubjectController) GetAll() {
	s.Data["json"] = BaseResponse{
		Code: 0,
		Msg:  msgSuccess,
		Data: models.Sm.GetAll(),
	}
	s.ServeJSON()
}
