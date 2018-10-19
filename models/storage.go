package models

import (
	"database/sql"
	"fmt"

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
	// "root:123456@tcp(localhost:3306)/lflss?charset=utf8"
	path := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
	logs.Debug(path)
	ma.db, err = sql.Open("mysql", path)
	if err != nil {
		panic("cannot connect to mysql")
	}
}

// load data
func (ma *mysqlAgent) LoadAllData() error {
	logs.Info("start loading data")
	// load all teachers
	teacherMap := make(map[UserID]*Teacher)
	{
		rows, err := ma.db.Query("SELECT iTeacherID,eGender,vName,vMobile FROM tbTeacher WHERE eStatus = 1;")
		if err != nil {
			logs.Warn("query database failed", "err", err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := &Teacher{}
			err = rows.Scan(&tmp.ID, &tmp.Gender, &tmp.LoginName, &tmp.Mobile)
			if err != nil {
				continue
			}
			teacherMap[tmp.ID] = tmp
		}
	}

	// init teacher manager
	Tm.Init(teacherMap)

	// load class
	classMap := make(map[ClassID]*Class)
	{
		rows, err := ma.db.Query("SELECT iClassID,iGrade,iIndex,vName FROM tbclass WHERE eStatus = 1;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := &Class{}
			err = rows.Scan(&tmp.ID, &tmp.Grade, &tmp.Index, &tmp.Name)
			if err != nil {
				continue
			}
			if tmp.ID != tmp.GetID() {
				logs.Warn("[LoadAllData] invalid id found", "id", tmp.ID, "expecting", tmp.GetID(), "grade", tmp.Grade, "index", tmp.Index)
				continue
			}
			classMap[tmp.ID] = tmp
		}
	}

	logs.Info("total class count is %d", len(classMap))

	// load class-teacher
	{
		rows, err := ma.db.Query("SELECT iClassID,iTeacherID FROM tbclassteacherrelation WHERE eStatus=1;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var classID ClassID
			var teacherID UserID
			err = rows.Scan(&classID, &teacherID)
			if err != nil {
				continue
			}
			// 检查teacherID 是否存在
			if _, ok := teacherMap[teacherID]; !ok {
				logs.Trace("data broken", "teacherID", teacherID)
				continue
			}

			// add to class map
			if v, ok := classMap[classID]; ok {
				v.TeacherIDs = append(v.TeacherIDs, teacherID)
			}
		}
	}

	// init class manager
	Cm.Init(classMap)

	// load students
	userMap := make(map[UserID]*User)
	{
		rows, err := ma.db.Query("SELECT iUserID, vUserName, vRegistNumber, eGender FROM tbuser WHERE eStatus = 1;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			u := User{}
			err = rows.Scan(&u.ID, &u.LoginName, &u.RegisterID, &u.Gender)
			if err != nil {
				continue
			}
			userMap[u.ID] = &u
		}
	}
	// init student manager
	Um.Init(userMap)

	// init access control
	userPassMap := make(map[string]string)
	teacherPassMap := make(map[string]string)
	{
		rows, err := ma.db.Query("SELECT iUserID, vPassword, eType FROM tbpassword;")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var id UserID
			var password string
			var eType int
			err = rows.Scan(&id, &password, &eType)
			if err != nil {
				logs.Error("scan failed", err)
				continue
			}
			if eType == 1 {
				if v, ok := userMap[id]; !ok {
					continue
				} else {
					userPassMap[v.LoginName] = password
				}
			} else if eType == 2 {
				if v, ok := teacherMap[id]; !ok {
					continue
				} else {
					teacherPassMap[v.LoginName] = password
				}
			}
		}
	}

	Ac.teacherMap = teacherPassMap
	Ac.studentMap = userPassMap
	logs.Info(userPassMap)
	logs.Info("load data success")
	return nil
}

// for teacher
func (ma *mysqlAgent) InsertTeacher(t *Teacher) error {
	// Prepare statement for inserting data
	stmtIns, err := ma.db.Prepare("INSERT INTO `tbteacher` (`eGender`, `vName`, `vMobile`) VALUES (?,?,?);")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(t.Gender, t.LoginName, t.Mobile)
	if err != nil {
		return err
	}
	id, err := resp.LastInsertId()
	t.ID = UserID(id)
	return nil
}

func (ma *mysqlAgent) UpdateTeacher(t *Teacher) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbteacher SET eGender=?,vName=?,vMobile=? WHERE iTeacherID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(t.Gender, t.LoginName, t.Mobile, t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ma *mysqlAgent) DeleteTeacher(teacherID UserID) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbteacher set eStatus=? WHERE iTeacherID=?;")
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
	stmtIns, err := ma.db.Prepare("INSERT INTO `tbClass` (`iClassID`, `iGrade`, `iIndex`, `vName`) VALUES (?,?,?,?);")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(t.ID, t.Grade, t.Index, t.Name)
	if err != nil {
		logs.Warn("execute sql failed", "err", err)
		return err
	}
	return nil
}

func (ma *mysqlAgent) UpdateClass(t *Class) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbClass SET vName=? WHERE iClassID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(t.Name, t.ID)
	if err != nil {
		logs.Warn("execute sql failed", "err", err)
		return err
	}
	return nil
}

func (ma *mysqlAgent) DeleteClass(id ClassID) error {
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
	return nil
}

// teacher-class relation
func (ma *mysqlAgent) AttachTeacher(id ClassID, teacherID []UserID) error {
	if ma == nil {
		return nil
	}

	if id == 0 || len(teacherID) == 0 {
		return ErrInvalidParam
	}

	stmtIns, err := ma.db.Prepare("INSERT INTO `tbClassTeacherRelation` (`iClassID`, `iTeacherID`) VALUES (?,?);")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	for _, v := range teacherID {
		_, err = stmtIns.Exec(id, v)
		if err != nil {
			logs.Warn("execute sql failed", "err", err)
			return err
		}
	}
	return nil
}

func (ma *mysqlAgent) DetachTeacher(id ClassID, teacherID []UserID) error {
	if ma == nil {
		return nil
	}

	if id == 0 || len(teacherID) == 0 {
		return ErrInvalidParam
	}

	stmtIns, err := ma.db.Prepare("UPDATE tbClassTeacherRelation set eStatus=? WHERE iClassID=? AND iTeacherID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	for _, v := range teacherID {
		_, err = stmtIns.Exec(id, v)
		if err != nil {
			logs.Warn("execute sql failed", "err", err)
			return err
		}
	}
	return nil
}

// student
// save student and password
func (ma *mysqlAgent) InsertUser(u *User) error {
	tx, err := ma.db.Begin()
	if err != nil {
		logs.Warn("[mysqlAgent::InsertUser] failed, err=", err)
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO `tbuser`(`vRegistNumber`, `vUserName`, `eGender`) VALUES (?,?,?)")
	if err != nil {
		logs.Error("[mysqlAgent::InsertUser] failed", "err")
		return err
	}

	rs, err := stmt.Exec(u.RegisterID, u.LoginName, u.Gender)
	if err != nil {
		logs.Warn("[mysqlAgent::InsertUser]failed", err)
		return err
	}

	id, err := rs.LastInsertId()
	u.ID = UserID(id)

	stmt, err = tx.Prepare("INSERT INTO `tbPassword`(`iUserID`, `vPassword`) VALUES (?,?)")
	if err != nil {
		logs.Warn("[mysqlAgent::InsertUser] failed", err)
		return err
	}

	rs, err = stmt.Exec(id, u.Password)
	if err != nil {
		logs.Warn("[mysqlAgent::InsertUser] failed, err=", err)
		return err
	}
	return tx.Commit()
}

func (ma *mysqlAgent) UpdateUser(u *User) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbUser SET vName=? WHERE iUserID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(u.LoginName, u.ID)
	if err != nil {
		logs.Warn("execute sql failed", "err", err)
		return err
	}
	return nil
}

func (ma *mysqlAgent) DeleteUser(uid UserID) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbUser SET eStatus=? WHERE iUserID=?;")
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

//
