package controllers

import (
	"github.com/astaxie/beego"
	"github.com/arong/dean/models"
	"encoding/json"
)

// Operations about object
type ClassController struct {
	beego.Controller
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Teacher	true		"The object content"
// @Success 200 {string} models.Teacher.ID
// @router / [post]
func (o *ClassController) Post() {
	request := models.Teacher{}
	resp := CommResp{Code: -1}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = invalidJSON
		goto Out
	}

	err = models.Tm.AddTeacher(&request)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = nil
Out:
	o.Data["json"] = resp
	o.ServeJSON()
}

