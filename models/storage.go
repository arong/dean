package models

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
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
			logs.Error("[LoadAllData] failed to load tbSubject", "err", err)
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
			logs.Error("[LoadAllData] failed to load tbTeacher", "err", err)
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
			logs.Error("[LoadAllData] failed to load tbClass", "err", err)
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
			logs.Error("[LoadAllData] failed to load tbClassTeacherRelation", "err", err)
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
		rows, err := ma.db.Query("SELECT iUserID,vName,vRegistNumber,eGender,iClassID FROM tbStudent WHERE eStatus = 1;")
		if err != nil {
			logs.Error("[LoadAllData] failed to load tbStudent", "err", err)
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
	loginMap := make(map[LoginKey]*LoginInfo)
	{
		rows, err := ma.db.Query("SELECT iUserID,eType,vLoginName,vPassword FROM tbPassword;")
		if err != nil {
			logs.Error("[LoadAllData] failed to load tbPassword", "err", err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := LoginInfo{}
			err = rows.Scan(&tmp.ID, &tmp.UserType, &tmp.LoginName, &tmp.Password)
			if err != nil {
				logs.Error("scan failed", err)
				continue
			}
			loginMap[LoginKey{UserType: tmp.UserType, LoginName: tmp.LoginName}] = &tmp
		}
	}
	Ac.loginMap = loginMap

	// init questionnaire
	questionMap := make(map[int]*QuestionnaireInfo)
	{
		rows, err := ma.db.Query("SELECT iQuestionnaireID,vTitle,dtStartTime,dtStopTime,eDraftStatus,vEditorName FROM tbQuestionnaire;")
		if err != nil {
			logs.Error("[LoadAllData] failed to load tbQuestionnaire", "err", err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := QuestionnaireInfo{}
			err = rows.Scan(&tmp.QuestionnaireID, &tmp.Title, &tmp.StartTime, &tmp.StopTime, &tmp.Status, &tmp.Editor)
			if err != nil {
				logs.Error("scan tbQuestionnaire failed", err)
				continue
			}
			decoded, err := base64.StdEncoding.DecodeString(tmp.Title)
			if err != nil {
				continue
			}
			tmp.Title = string(decoded)
			questionMap[tmp.QuestionnaireID] = &tmp
		}
	}
	QuestionnaireManager.Init(questionMap)

	// init question
	{
		rows, err := ma.db.Query("SELECT iQuestionID, iQuestionnaireID, iIndex, eType, bRequired, vQuestion, vContent FROM tbQuestion;")
		if err != nil {
			logs.Error("[LoadAllData] failed to load tbQuestion", "err", err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := QuestionInfo{}
			buff := ""
			require := 0
			err = rows.Scan(&tmp.QuestionID, &tmp.QuestionnaireID, &tmp.Index, &tmp.Type, &require, &tmp.Question, &buff)
			if err != nil {
				logs.Error("[LoadAllData] scan tbQuestion failed", err)
				continue
			}
			if require == 1 {
				tmp.Required = true
			}

			decoded, err := base64.StdEncoding.DecodeString(tmp.Question)
			if err != nil {
				logs.Warn("[LoadAllData] tbQuestion.Question data error", "err", err)
				continue
			}
			tmp.Question = string(decoded)

			decoded, err = base64.StdEncoding.DecodeString(buff)
			if err != nil {
				logs.Warn("[LoadAllData] tbQuestion.Options data error", "err", err)
				continue
			}
			err = json.Unmarshal(decoded, &tmp.Options)
			if err != nil {
				logs.Warn("[LoadAllData] invalid Option data", "err", err)
				continue
			}

			if q, ok := QuestionnaireManager.questionnaires[tmp.QuestionnaireID]; ok {
				q.Questions = append(q.Questions, &tmp)
			} else {
				logs.Warn("[LoadAllData] questionnaire id not found")
				continue
			}
			QuestionnaireManager.questions[tmp.QuestionID] = &tmp
		}
	}
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
func (ma *mysqlAgent) InsertPassword(l *LoginInfo) error {
	stmtIns, err := ma.db.Prepare("INSERT INTO `tbPassword`(`iUserID`, `eType`, `vLoginName`, `vPassword`) VALUES (?,?,?,?);")
	if err != nil {
		logs.Warn("[InsertPassword] Prepare sql failed", "err", err)
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(l.ID, l.UserType, l.LoginName, l.Password)
	if err != nil {
		logs.Warn("[InsertPassword] execute sql failed", "err", err)
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

func (ma *mysqlAgent) ResetAllPassword(password string) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbPassword SET vPassword=? WHERE eType=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(password, base.AccountTypeStudent)
	if err != nil {
		logs.Warn("[DropAllPassword] execute sql failed", "err", err)
		return err
	}

	count, err := resp.RowsAffected()
	if err != nil {
		logs.Warn("[DropAllPassword] database failed", "err", err)
	}
	logs.Info("[DropAllPassword] rows affected", "count", count)
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

// Questionnaire zone

// AddQuestionnaire add new questionnaire to current system
func (ma *mysqlAgent) AddQuestionnaire(q *QuestionnaireInfo) error {

	stmtIns, err := ma.db.Prepare("INSERT INTO tbQuestionnaire (`vTitle`,`dtStartTime`,`dtStopTime`,`eDraftStatus`,`vEditorName`) VALUES (?,?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(base64.StdEncoding.EncodeToString([]byte(q.Title)), q.StartTime, q.StopTime, q.Status, q.Editor)
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
	return nil
}

// UpdateQuestionnaire modify questionnaire info, not its
func (ma *mysqlAgent) UpdateQuestionnaire(q *QuestionnaireInfo) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbQuestionnaire SET `vTitle`=?,`dtStartTime`=?,`dtStopTime`=?,`vEditorName`=? WHERE iQuestionnaireID=? AND `eDraftStatus`=1")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(base64.StdEncoding.EncodeToString([]byte(q.Title)), q.StartTime, q.StopTime, q.Editor, q.QuestionnaireID)
	if err != nil {
		logs.Warn("[UpdateQuestionnaire] execute sql failed", "err", err)
		return err
	}

	rowAffected, err := resp.RowsAffected()
	if err != nil {
		logs.Warn("[UpdateQuestionnaire] database error")
		return err
	}

	if rowAffected != 1 {
		logs.Warn("[UpdateQuestionnaire] data error")
	}
	return nil
}

func (ma *mysqlAgent) ExpireQuestionnaire(q *QuestionnaireInfo) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbQuestionnaire (`vTitle`,`dtStartTime`,`dtStopTime`,`eDraftStatus`) VALUES (?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

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
	return nil
}

func (ma *mysqlAgent) DeleteQuestionnaire(id int) error {
	stmtIns, err := ma.db.Prepare("DELETE FROM tbQuestionnaire WHERE `iQuestionnaireID`=? AND `eDraftStatus`=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(id, QStatusDraft)
	if err != nil {
		logs.Warn("[DeleteQuestionnaire] execute sql failed", "err", err)
		return err
	}

	count, err := resp.RowsAffected()
	if err != nil {
		logs.Warn("[DeleteQuestionnaire] unexpected update rows", "err", err)
		return err
	}

	if count != 1 {
		logs.Warn("[DeleteQuestionnaire] duplicated data found")
		return errors.New("data error")
	}
	return nil
}

func (ma *mysqlAgent) AddQuestion(questionnaireID int, info *QuestionInfo) (int, error) {
	buff, err := json.Marshal(info.Options)
	if err != nil {
		return 0, err
	}
	encoded := base64.StdEncoding.EncodeToString(buff)

	stmtIns, err := ma.db.Prepare("INSERT INTO tbQuestion (`iQuestionnaireID`,`vQuestion`,`iIndex`,`eType`,`bRequired`,`vContent`) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(questionnaireID, base64.StdEncoding.EncodeToString([]byte(info.Question)), info.Index, info.Type, func() int {
		if info.Required {
			return 1
		} else {
			return 0
		}
	}(), encoded)

	if err != nil {
		logs.Warn("[AddQuestion] execute sql failed", "err", err)
		return 0, err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		logs.Warn("[AddQuestion] unexpected update rows", "err", err)
		return 0, err
	}

	return int(id), nil
}

func (ma *mysqlAgent) UpdateQuestion(info *QuestionInfo) (int, error) {
	buff, err := json.Marshal(info.Options)
	if err != nil {
		return 0, err
	}
	encoded := base64.StdEncoding.EncodeToString(buff)
	stmtIns, err := ma.db.Prepare("UPDATE tbQuestion SET `vQuestion`=?,`iIndex`=?,`eType`=?,`bRequired`=?,`vContent`=? WHERE iQuestionID=?")
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(base64.StdEncoding.EncodeToString([]byte(info.Question)), info.Index, info.Type, func() int {
		if info.Required {
			return 1
		} else {
			return 0
		}
	}(), encoded, info.QuestionID)

	if err != nil {
		logs.Warn("[UpdateQuestion] execute sql failed", "err", err)
		return 0, err
	}

	id, err := resp.LastInsertId()
	if err != nil {
		logs.Warn("[UpdateQuestion] unexpected update rows", "err", err)
		return 0, err
	}

	return int(id), nil
}

func (ma *mysqlAgent) DeleteQuestion(questionID int) error {
	stmtIns, err := ma.db.Prepare("DELETE FROM tbQuestion WHERE iQuestionID = ?")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(questionID)

	if err != nil {
		logs.Warn("[DeleteQuestion] execute sql failed", "err", err)
		return err
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		logs.Warn("[DeleteQuestion] unexpected update rows", "err", err)
		return err
	}

	if rows != 1 {
		logs.Warn("[DeleteQuestion] data error", "rows", rows)
	}
	return nil
}
