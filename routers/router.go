// @APIVersion 1.0.0
// @Title Aronic dean system for LFLSS
// @Description Dean is designed for the managing dean issue in lflss
// @Contact aronic@outlook.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/arong/dean/controllers"
	"github.com/astaxie/beego"
)

const (
	// fixedPath is the top directory for all the resource
	fixedPath = "/api/v1/dean"
)

func init() {
	ns := beego.NewNamespace(fixedPath,
		beego.NSNamespace("/class",
			beego.NSInclude(
				&controllers.ClassController{},
			),
		),
		beego.NSNamespace("/vote",
			beego.NSInclude(
				&controllers.VoteController{},
			),
		),
		beego.NSNamespace("/teacher",
			beego.NSInclude(
				&controllers.TeacherController{},
			),
		),
		beego.NSNamespace("/score",
			beego.NSInclude(
				&controllers.ScoreController{},
			),
		),
		beego.NSNamespace("/student",
			beego.NSNamespace("/score",
				beego.NSInclude(
					&controllers.StudentScoreController{},
				),
			),
			beego.NSInclude(
				&controllers.StudentController{},
			),
		),
		beego.NSNamespace("/subject",
			beego.NSInclude(
				&controllers.SubjectController{},
			),
		),
		beego.NSNamespace("/questionnaire",
			beego.NSNamespace("/question",
				beego.NSInclude(
					&controllers.QuestionController{},
				),
			),
			beego.NSNamespace("/view",
				beego.NSInclude(
					&controllers.ScoreController{},
				),
			),
			beego.NSNamespace("/edit",
				beego.NSInclude(
					&controllers.QuestionnaireController{},
				),
			),
			beego.NSNamespace("/vote",
				beego.NSInclude(
					&controllers.VoteController{},
				),
			),
		),
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
