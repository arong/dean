package models

import (
	"database/sql"
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

func (ma *mysqlAgent) Init() {
	var err error
	ma.db, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/lflss?charset=utf8")
	if err != nil {
		panic("cannot connect to mysql")
	}
}

func (ma *mysqlAgent) loadAllTeachers() ([]*Teacher, error) {
	ret := []*Teacher{}
	rows, err := ma.db.Query("SELECT iTeacherID,eGender,vName,vMobile FROM tbTeacher WHERE eStatus = 1;")
	if err != nil {
		return ret, err
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &Teacher{}
		err = rows.Scan(&tmp.ID, &tmp.Gender, &tmp.Name, &tmp.Mobile)
		if err != nil {
			continue
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

func (ma *mysqlAgent) InsertTeacher(t *Teacher) error {
	// Prepare statement for inserting data
	stmtIns, err := ma.db.Prepare("INSERT INTO `tbteacher` (`eGender`, `vName`, `vMobile`) VALUES (?,?,?);")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	resp, err := stmtIns.Exec(t.ID, t.Gender, t.Name, t.Mobile)
	if err != nil {
		return err
	}
	t.ID, err = resp.LastInsertId()
	return nil
}

func (ma *mysqlAgent) UpdateTeacher(t *Teacher) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbteacher SET eGender=?,vName=?,vMobile=? WHERE iTeacherID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(t.Gender, t.Name, t.Mobile, t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ma *mysqlAgent) DeleteTeacher(t *Teacher) error {
	stmtIns, err := ma.db.Prepare("UPDATE tbteacher set eStatus=? WHERE iTeacherID=?;")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(eStatusDeleted, t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ma *mysqlAgent) loadAllClass() (map[ClassID]*Class, error) {
	ret := make(map[ClassID]*Class)

	{

		rows, err := ma.db.Query("SELECT iClassID,iGradeNumber,iClassNumber,vClassName FROM tbclass WHERE eStatus = 1;")
		if err != nil {
			return ret, err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := &Class{}
			err = rows.Scan(&tmp.ID, &tmp.Grade, &tmp.Index, &tmp.Name)
			if err != nil {
				continue
			}
			ret[tmp.ID] = tmp
			//ret = append(ret, tmp)
		}
	}

	{

		rows, err := ma.db.Query("SELECT iClassID,iTeacherID FROM tbclassteacherrelation WHERE eStatus=1;")
		if err != nil {
			return ret, err
		}
		defer rows.Close()

		for rows.Next() {
			var classID ClassID
			var teacherID int64
			err = rows.Scan(&classID, &teacherID)
			if err != nil {
				continue
			}
			if v, ok := ret[classID]; ok {
				v.TeacherIDs = append(v.TeacherIDs, teacherID)
			}
		}
	}

	return ret, nil
}