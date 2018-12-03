package models

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

var Ma mysqlAgent

const (
	eStatusDeleted = 2
	eStatusNormal  = 1
)

type mysqlAgent struct {
	db *sql.DB
}

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"DBName"`
}

func (c *DBConfig) GetConf() error {
	c.User = beego.AppConfig.String("user")
	c.Password = beego.AppConfig.String("password")
	c.Host = beego.AppConfig.String("host")
	c.Port, _ = beego.AppConfig.Int("port")
	c.DBName = beego.AppConfig.String("dbName")

	logs.Debug(*c)
	return nil
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

// load data
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
	teacherMap := make(map[UserID]*Teacher)
	{
		rows, err := ma.db.Query("SELECT iTeacherID,eGender,vName,vMobile,iPrimarySubjectID,dtBirthday,vAddress FROM tbTeacher WHERE eStatus = 1;")
		if err != nil {
			logs.Warn("query database failed", "err", err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := &Teacher{}
			err = rows.Scan(&tmp.TeacherID, &tmp.Gender, &tmp.RealName, &tmp.Mobile, &tmp.SubjectID, &tmp.Birthday, &tmp.Address)
			if err != nil {
				continue
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
		rows, err := ma.db.Query("SELECT iUserID, vUserName, vRegistNumber, eGender FROM tbStudent WHERE eStatus = 1;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			u := StudentInfo{}
			err = rows.Scan(&u.StudentID, &u.RealName, &u.RegisterID, &u.Gender)
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

	logs.Info("load data success")
	return nil
}

// for teacher
func (ma *mysqlAgent) InsertTeacher(t *Teacher) error {
	// Prepare statement for inserting data
	stmtIns, err := ma.db.Prepare("INSERT INTO `tbTeacher` (`eGender`, `vName`, `vMobile`, `dtBirthday`, `vAddress`, `iPrimarySubjectID`) VALUES (?,?,?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(t.Gender, t.RealName, t.Mobile, t.Birthday, t.Address, t.SubjectID)
	if err != nil {
		return err
	}
	id, err := resp.LastInsertId()
	t.TeacherID = UserID(id)
	return nil
}

func (ma *mysqlAgent) UpdateTeacher(t *Teacher) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbTeacher SET eGender=?,vName=?,vMobile=?,dtBirthday=?,iPrimarySubjectID=?, vAddress=? WHERE iTeacherID=?;")
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

func (ma *mysqlAgent) DeleteTeacher(teacherID UserID) error {
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

// class
// 增删改查
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

// student
// save student and password
func (ma *mysqlAgent) InsertUser(u *StudentInfo) (int64, error) {
	stmt, err := ma.db.Prepare("INSERT INTO `tbStudent`(`vRegistNumber`, `vUserName`, `eGender`) VALUES (?,?,?)")
	if err != nil {
		logs.Error("[mysqlAgent::InsertUser] failed", "err")
		return 0, err
	}

	rs, err := stmt.Exec(u.RegisterID, u.RealName, u.Gender)
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

func (ma *mysqlAgent) UpdateUser(u *StudentInfo) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbStudent SET vName=? WHERE iUserID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(u.RealName, u.StudentID)
	if err != nil {
		logs.Warn("execute sql failed", "err", err)
		return err
	}
	return nil
}

func (ma *mysqlAgent) DeleteUser(uid int64) error {
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

// password
func (ma *mysqlAgent) UpdatePassword(id UserID, password string) error {
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
