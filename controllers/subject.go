package controllers

import (
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
)

// Operations about Users
type SubjectController struct {
	beego.Controller
}

// @Title GetAll
// @Description get all Users
// @Param	grade		query 	string	true		"The grade of class"
// @Param	index		query 	string	true		"The number of class"
// @Success 200 {object} models.User
// @router / [get]
func (s *SubjectController) GetAll() {
	s.Data["json"] = CommResp{
		Code: 0,
		Msg:  msgSuccess,
		Data: models.Sm.GetAll(),
	}
	s.ServeJSON()
}
