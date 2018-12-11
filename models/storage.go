package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/arong/dean/base"
	"github.com/bearbin/go-age"
	"github.com/pkg/errors"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

// Ma mysql agent
var Ma mysqlAgent

const (
	eStatusDeleted = 2
	eStatusNormal  = 1
)

type mysqlAgent struct {
	db *sql.DB
}

func (ma *mysqlAgent) Init(conf *DBConfig) {
	var err error
	// example: "root:123456@tcp(localhost:3306)/lflss?charset=utf8"
	path := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
	ma.db, err = sql.Open("mysql", path)
	if err != nil {
		panic("cannot connect to mysql")
	}
}

// LoadAllData load data
func (ma *mysqlAgent) LoadAllData() error {
	logs.Info("start loading data")

	// load all subject info
	subjectMap := make(map[int]SubjectInfo)
	{
		rows, err := ma.db.Query("select iSubjectID,vSubjectKey, vSubjectName from tbSubject where eStatus =1;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := SubjectInfo{}
			err = rows.Scan(&tmp.ID, &tmp.Key, &tmp.Name)
			if err != nil {
				logs.Error("scan failed", err)
				continue
			}
			subjectMap[tmp.ID] = tmp
		}
	}
	Sm.subject = subjectMap

	// load all teachers
	teacherMap := make(map[int64]*Teacher)
	{
		rows, err := ma.db.Query("SELECT iTeacherID,eGender,vName,vMobile,iSubjectID,dtBirthday,vAddress FROM tbTeacher WHERE eStatus = 1;")
		if err != nil {
			logs.Warn("query database failed", "err", err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := &Teacher{}
			err = rows.Scan(&tmp.TeacherID, &tmp.Gender, &tmp.RealName, &tmp.Mobile, &tmp.SubjectID, &tmp.Birthday, &tmp.Address)
			if err != nil {
				logs.Warn("[LoadAllData] data error at tbTeacher")
				continue
			}
			if tmp.Birthday != "" {
				birth, err := time.Parse(base.DateFormat, tmp.Birthday)
				if err != nil {
					logs.Warn("[LoadAllData] data error at tbTeacher", "birthday", tmp.Birthday, "err", err)
					continue
				}
				tmp.Age = age.Age(birth)
			}
			teacherMap[tmp.TeacherID] = tmp
		}
	}

	// init teacher manager
	Tm.Init(teacherMap)

	// load class
	classMap := make(map[int]*Class)
	{
		rows, err := ma.db.Query("SELECT iClassID,iGrade,iIndex,vName,iMasterID,iStartYear,eTerm FROM tbClass WHERE eStatus = 1;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := &Class{}
			err = rows.Scan(&tmp.ID, &tmp.Grade, &tmp.Index, &tmp.Name, &tmp.MasterID, &tmp.Year, &tmp.Term)
			if err != nil {
				continue
			}
			classMap[tmp.ID] = tmp
		}
	}
	logs.Info("total class count is %d", len(classMap))

	// load class-teacher
	{
		rows, err := ma.db.Query("SELECT iClassID,iTeacherID,iSubjectID FROM tbClassTeacherRelation WHERE eStatus=1;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var classID int
			tmp := InstructorInfo{}
			err = rows.Scan(&classID, &tmp.TeacherID, &tmp.SubjectID)
			if err != nil {
				continue
			}
			// 检查teacherID 是否存在
			if _, ok := teacherMap[tmp.TeacherID]; !ok {
				logs.Warn("data broken, teacher not found", "teacherID", tmp.TeacherID)
				continue
			}

			// 检查classID 是否存在
			if v, ok := classMap[classID]; !ok {
				logs.Warn("data broken, class not found", "classID", classID)
				continue
			} else {
				v.TeacherList = append(v.TeacherList, tmp)
			}
		}
	}

	// init class manager
	Cm.Init(classMap)

	// load students
	userMap := make(map[int64]*StudentInfo)
	{
		rows, err := ma.db.Query("SELECT iUserID,vUserName,vRegistNumber,eGender,iClassID FROM tbStudent WHERE eStatus = 1;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			u := StudentInfo{}
			err = rows.Scan(&u.StudentID, &u.RealName, &u.RegisterID, &u.Gender, &u.ClassID)
			if err != nil {
				continue
			}
			userMap[u.StudentID] = &u
		}
	}
	// init student manager
	Um.Init(userMap)

	// init access control
	loginMap := make(map[string]*loginInfo)
	{
		rows, err := ma.db.Query("SELECT iUserID,eType,vLoginName,vPassword FROM tbPassword;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := loginInfo{}
			err = rows.Scan(&tmp.ID, &tmp.UserType, &tmp.LoginName, &tmp.Password)
			if err != nil {
				logs.Error("scan failed", err)
				continue
			}
			loginMap[tmp.LoginName] = &tmp
		}
	}
	Ac.loginMap = loginMap

	// init questionnaire
	questionMap := make(map[int]*QuestionnaireInfo)
	{
		rows, err := ma.db.Query("SELECT iQuestionnaireID,vTitle,dtStartTime,dtStopTime,eDraftStatus FROM tbQuestionnaire;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := QuestionnaireInfo{}
			err = rows.Scan(&tmp.QuestionnaireID, &tmp.Title, &tmp.StartTime, &tmp.StopTime, &tmp.Status)
			if err != nil {
				logs.Error("scan failed", err)
				continue
			}
			questionMap[tmp.QuestionnaireID] = &tmp
		}
	}
	Qm.questionnaires = questionMap

	logs.Info("load data success")
	return nil
}

// InsertTeacher insert teacher info
func (ma *mysqlAgent) InsertTeacher(t *Teacher) error {
	// Prepare statement for inserting data
	stmtIns, err := ma.db.Prepare("INSERT INTO `tbTeacher` (`eGender`, `vName`, `vMobile`, `dtBirthday`, `vAddress`, `iSubjectID`) VALUES (?,?,?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(t.Gender, t.RealName, t.Mobile, t.Birthday, t.Address, t.SubjectID)
	if err != nil {
		return err
	}
	id, err := resp.LastInsertId()
	t.TeacherID = id
	return nil
}

// UpdateTeacher update teacher info
func (ma *mysqlAgent) UpdateTeacher(t *Teacher) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbTeacher SET eGender=?,vName=?,vMobile=?,dtBirthday=?,iSubjectID=?, vAddress=? WHERE iTeacherID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(t.Gender, t.RealName, t.Mobile, t.Birthday, t.SubjectID, t.Address, t.TeacherID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTeacher delete teacher info
func (ma *mysqlAgent) DeleteTeacher(teacherID int64) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbTeacher set eStatus=? WHERE iTeacherID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(eStatusDeleted, teacherID)
	if err != nil {
		return err
	}
	return nil
}

// InsertClass insert class info
func (ma *mysqlAgent) InsertClass(t *Class) error {
	// Prepare statement for inserting data
	stmtIns, err := ma.db.Prepare("INSERT INTO `tbClass` (`iGrade`, `iIndex`, `vName`,`iMasterID`,`iStartYear`,`eTerm`) VALUES (?,?,?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(t.Grade, t.Index, t.Name, t.MasterID, t.Year, t.Term)
	if err != nil {
		logs.Warn("execute sql failed", "err", err)
		return err
	}

	id, err := resp.LastInsertId()
	t.ID = int(id)
	// insert into class teacher relation table
	{
		stmtIns, err := ma.db.Prepare("INSERT INTO `tbClassTeacherRelation` (`iClassID`,`iSubjectID`, `iTeacherID`) VALUES (?,?,?);")
		if err != nil {
			return err
		}

		for _, v := range t.TeacherList {
			_, err = stmtIns.Exec(t.ID, v.SubjectID, v.TeacherID)
			if err != nil {
				logs.Warn("execute sql failed", "err", err)
				return err
			}
		}
	}
	return nil
}

// UpdateClass update class info
func (ma *mysqlAgent) UpdateClass(t *Class) error {
	// insert into tbClass
	{
		stmtIns, err := ma.db.Prepare("UPDATE tbClass SET vName=?,iMasterID=?,iStartYear=?,eTerm=? WHERE iClassID=?;")
		if err != nil {
			return err
		}
		defer stmtIns.Close()

		_, err = stmtIns.Exec(t.Name, t.MasterID, t.Year, t.Term, t.ID)
		if err != nil {
			logs.Warn("execute sql failed", "err", err)
			return err
		}
	}

	// remove existing item
	{
		stmtIns, err := ma.db.Prepare("DELETE FROM `tbClassTeacherRelation` WHERE `iClassID`=? AND `iSubjectID`=? AND `iTeacherID`=?;")
		if err != nil {
			return err
		}

		for _, v := range t.RemoveList {
			_, err = stmtIns.Exec(t.ID, v.SubjectID, v.TeacherID)
			if err != nil {
				logs.Warn("execute sql failed", "err", err)
				return err
			}
		}
	}

	// insert into class teacher relation table
	{
		stmtIns, err := ma.db.Prepare("INSERT INTO `tbClassTeacherRelation` (`iClassID`,`iSubjectID`, `iTeacherID`) VALUES (?,?,?);")
		if err != nil {
			return err
		}

		for _, v := range t.AddList {
			_, err = stmtIns.Exec(t.ID, v.SubjectID, v.TeacherID)
			if err != nil {
				logs.Warn("execute sql failed", "err", err)
				return err
			}
		}
	}
	return nil
}

// DeleteClass delete class info
func (ma *mysqlAgent) DeleteClass(id int) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbClass set eStatus=? WHERE iClassID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(eStatusDeleted, id)
	if err != nil {
		logs.Warn("execute sql failed", "err", err)
		return err
	}

	// remove teacher relation
	{
		stmtIns, err := ma.db.Prepare("DELETE FROM tbClassTeacherRelation WHERE iClassID=?;")
		if err != nil {
			return err
		}
		defer stmtIns.Close()

		_, err = stmtIns.Exec(id)
		if err != nil {
			logs.Warn("execute sql failed", "err", err)
			return err
		}
	}
	return nil
}

// InsertStudent insert teacher info
func (ma *mysqlAgent) InsertStudent(u *StudentInfo) (int64, error) {
	stmt, err := ma.db.Prepare("INSERT INTO `tbStudent`(`vRegistNumber`, `vName`, `eGender`,`iClassID`,`vAddress`,`dtBirthday`) VALUES (?,?,?,?,?,?)")
	if err != nil {
		logs.Error("[mysqlAgent::InsertUser] failed", "err")
		return 0, err
	}

	rs, err := stmt.Exec(u.RegisterID, u.RealName, u.Gender, u.ClassID, u.Address, u.Birthday)
	if err != nil {
		logs.Warn("[mysqlAgent::InsertUser]failed", err)
		return 0, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		logs.Warn("[mysqlAgent::InsertUser] LastInsertId failed", err)
		return 0, err
	}
	return id, nil
}

// UpdateStudent update student info
func (ma *mysqlAgent) UpdateStudent(u *StudentInfo) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbStudent SET vRegistNumber=?,vName=?,eGender=?,iClassID=?,vAddress=?,dtBirthday=? WHERE iUserID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(u.RegisterID, u.RealName, u.Gender, u.ClassID, u.Address, u.Birthday, u.StudentID)
	if err != nil {
		logs.Warn("execute sql failed", "err", err)
		return err
	}
	return nil
}

// DeleteUser delete student info
func (ma *mysqlAgent) DeleteStudent(uid int64) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbStudent SET eStatus=? WHERE iUserID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(eStatusDeleted, uid)
	if err != nil {
		logs.Warn("execute sql failed", "err", err)
		return err
	}
	return nil
}

// UpdatePassword password
func (ma *mysqlAgent) UpdatePassword(id int64, password string) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbPassword SET vPassword=? WHERE iUserID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(password, id)
	if err != nil {
		logs.Warn("execute sql failed", "err", err)
		return err
	}
	return nil
}

// AddSubject add subject info
func (ma *mysqlAgent) AddSubject(s *SubjectInfo) (int, error) {
	stmtIns, err := ma.db.Prepare("INSERT INTO tbSubject (`vSubjectKey`, `vSubjectName`) VALUES (?,?);")
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(s.Key, s.Name)
	if err != nil {
		logs.Warn("execute sql failed", "err", err)
		return 0, err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		logs.Info("[] unexpected error", "err", err)
		return 0, nil
	}
	return int(id), nil
}

// UpdateSubject update subject info
func (ma *mysqlAgent) UpdateSubject(s *SubjectInfo) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbSubject SET `vSubjectKey`=? WHERE `iSubjectID`=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(s.Key, s.ID)
	if err != nil {
		logs.Warn("[mysqlAgent::UpdateSubject] execute sql failed", "err", err)
		return err
	}

	count, err := resp.RowsAffected()
	if err != nil {
		logs.Warn("[mysqlAgent::UpdateSubject] unexpected update rows", "err", err)
		return err
	}

	if count != 1 {
		logs.Warn("[mysqlAgent::UpdateSubject] duplicated data found")
		return errors.New("data error")
	}
	return nil
}

// DeleteSubject delete subject info
func (ma *mysqlAgent) DeleteSubject(id int) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbSubject SET `eStatus`=? WHERE `iSubjectID`=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(eStatusDeleted, id)
	if err != nil {
		logs.Warn("[DeleteSubject] execute sql failed", "err", err)
		return err
	}

	count, err := resp.RowsAffected()
	if err != nil {
		logs.Warn("[DeleteSubject] unexpected update rows", "err", err)
		return err
	}

	if count != 1 {
		logs.Warn("[DeleteSubject] duplicated data found")
		return errors.New("data error")
	}
	return nil
}

func (ma *mysqlAgent) AddQuestionnaire(q *QuestionnaireInfo) error {
	stmtIns, err := ma.db.Prepare("INSERT tbQuestionnaire (`vTitle`,`dtStartTime`,`dtStopTime`,`eDraftStatus`) VALUES (?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	insQuestion, err := ma.db.Prepare("INSERT INTO `tbQuestion` (`iQuestionnaireID`,`iIndex`,`vQuestion`, `eType`,`bRequired`) VALUES (?,?,?,?);")
	if err != nil {
		return err
	}

	insOption, err := ma.db.Prepare("INSERT INTO `tbOption` (`iQuestionID`,`iIndex`,`vOption`) VALUES (?,?,?);")
	if err != nil {
		return err
	}
	resp, err := stmtIns.Exec(eStatusDeleted, q.Title, q.StartTime, q.StopTime, q.Status)
	if err != nil {
		logs.Warn("[AddQuestionnaire] execute sql failed", "err", err)
		return err
	}

	tmp, err := resp.LastInsertId()
	if err != nil {
		logs.Warn("[AddQuestionnaire] insert tbQuestionnaire error", "err", err)
		return err
	}
	q.QuestionnaireID = int(tmp)

	// insert into class teacher relation table
	{
		for _, v := range q.Questions {
			resp, err = insQuestion.Exec(q.QuestionnaireID, v.Index, v.Question, v.Type, v.Required)
			if err != nil {
				logs.Warn("[AddQuestionnaire] execute sql failed", "err", err)
				return err
			}
			tmp, err = resp.LastInsertId()
			if err != nil {
				logs.Error("[AddQuestionnaire] insert tbQuestion failed", "question", v)
				return err
			}
			v.QuestionID = int(tmp)

			for _, val := range v.Options {
				resp, err := insOption.Exec(v.QuestionID, val.Index, val.Option)
				if err != nil {
					logs.Error("[AddQuestionnaire] insert tbOption error")
					return err
				}

				tmp, err = resp.LastInsertId()
				if err != nil {
					logs.Error("[AddQuestionnaire] insert tbOption error")
					return err
				}

				val.OptionID = int(tmp)
			}
		}
	}
	return nil
}
