package controllers

import (
	"encoding/json"
	"github.com/arong/dean/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
)

// Operations about object
type TeacherController struct {
	beego.Controller
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Teacher	true		"The object content"
// @Success 200 {string} models.Teacher.TeacherID
// @router / [post]
func (o *TeacherController) Post() {
	request := models.Teacher{}
	resp := CommResp{Code: -1}

	err := json.Unmarshal(o.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[TeacherController] Unmarshal failed", "err", err)
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

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.Teacher	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *TeacherController) Put() {
	resp := &CommResp{Code: -1}
	tmp := u.GetString(":uid")
	var teacher models.Teacher
	uid, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		logs.Debug("[TeacherController::Put] parse uid failed")
		goto Out
	}

	if uid == 0 {
		logs.Info("[TeacherController::Put] invalid teacher id")
		goto Out
	}

	err = json.Unmarshal(u.Ctx.Input.RequestBody, &teacher)
	if err != nil {
		logs.Info("[TeacherController::Put] unmarshal failed", "err", err)
		resp.Msg = msgInvalidJSON
		goto Out
	}

	teacher.TeacherID = models.UserID(uid)
	err = models.Tm.ModTeacher(&teacher)
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title Get
// @Description find object by objectid
// @Param	teacherID		path 	string	true		"the teacherID you want to get"
// @Success 200 {object}	models.Teacher
// @Failure 403 :teacherID is empty
// @router /:teacherID [get]
func (o *TeacherController) Get() {
	resp := CommResp{Code: -1}
	var err error
	var id int64
	ret := &models.Teacher{}

	teacherID := o.Ctx.Input.Param(":teacherID")
	if teacherID == "" {
		resp.Msg = msgInvalidParam
		goto Out
	}

	id, err = strconv.ParseInt(teacherID, 10, 64)
	if err != nil {
		resp.Msg = msgInvalidParam
		goto Out
	}

	ret, err = models.Tm.GetTeacherInfo(models.UserID(id))
	if err != nil {
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = ret
Out:
	o.Data["json"] = resp
	o.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Param	page		query 	string	true		"The username for login"
// @Param	size		query 	string	true		"The password for login"
// @Param	body		body 	models.TeacherFilter	true		"The object content"
// @Success 200 {object} models.TeacherListResp
// @router /filter [post]
func (o *TeacherController) GetAll() {
	resp := &CommResp{}
	request := models.TeacherFilter{}
	err := json.Unmarshal(o.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[TeacherController::GetAll] Unmarshal failed", "err", err)
		goto Out
	}

	if request.Page <= 0 || request.Size <= 0 {
		resp.Msg = msgInvalidParam
		logs.Debug("[TeacherController::GetAll] invalid size")
		goto Out
	}

	resp.Data = models.Tm.GetAll(&request)

Out:
	o.Data["json"] = resp
	o.ServeJSON()
}

type DeleteTeacherReq struct {
	IDList []models.UserID `json:"id_list"`
}

type DeleteTeacherResp struct {
	FailedList []models.UserID
}

// @Title Delete
// @Description delete the user
// @Param	body		body 	models.TeacherFilter	true		"The object content"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /delete [post]
func (tc *TeacherController) Delete() {
	request:= DeleteTeacherReq{}
	resp := &CommResp{Code: -1}
	ret := DeleteTeacherResp{}

	err := json.Unmarshal(tc.Ctx.Input.RequestBody, &request)
	if err != nil {
		resp.Msg = msgInvalidJSON
		logs.Debug("[TeacherController::Delete] Unmarshal failed", "err", err)
		goto Out
	}

	ret.FailedList, err = models.Tm.DelTeacher(request.IDList)
	if err != nil {
		logs.Debug("[TeacherController::Delete] failed", "err", err)
		resp.Msg = err.Error()
		resp.Data = ret
		goto Out
	}

	resp.Code = 0
	resp.Msg = msgSuccess
Out:
	tc.Data["json"] = resp
	tc.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *TeacherController) Login() {
	resp := &CommResp{Code: -1}
	username := u.GetString("username")
	password := u.GetString("password")

	token, err := models.Ac.Login(username, password, models.TypeTeacher)
	if err != nil {
		logs.Debug("[UserController::Login] login failed", username, err)
		resp.Msg = err.Error()
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess
	resp.Data = token
Out:
	u.Data["json"] = resp
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Param	token		query 	string	true		"The username for login"
// @Success 200 {string} logout success
// @router /logout [get]
func (u *TeacherController) Logout() {
	resp := &CommResp{Code: -1}
	token := u.GetString("username")
	if token == "" {
		logs.Debug("no token")
		resp.Msg = "invalid token"
		goto Out
	}

	if models.Ac.Logout(token) != nil {
		logs.Debug("logout failed")
		goto Out
	}
	resp.Code = 0
	resp.Msg = msgSuccess

Out:
	u.Data["json"] = resp
	u.ServeJSON()
}
