package controllers

import (
	"github.com/astaxie/beego"
	"github.com/arong/dean/models"
	"strconv"
	"fmt"
)

// Operations about object
type ScoreController struct {
	beego.Controller
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.ScoreInfo
// @Failure 403 :objectId is empty
// @router /:teacherID [get]
func (s *ScoreController) Get() {
	resp := CommResp{Code: -1}
	var err error
	var id int64
	ret := &models.ScoreInfo{}
	teacherID := s.Ctx.Input.Param(":teacherID")
	if teacherID == "" {
		resp.Msg = invalidParam
		goto Out
	}
	id, err = strconv.ParseInt(teacherID, 10, 64)
	if err != nil {
		resp.Msg = invalidParam
		goto Out
	}
	fmt.Println("teacherID=", id)
	ret, err = models.Vm.GetScore(id)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = ret
Out:
	s.Data["json"] = resp
	s.ServeJSON()
}
