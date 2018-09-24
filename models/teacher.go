package models

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

// a global handler
var Tm TeacherManager
var db *sql.DB

func init() {
	var err error
	//orm.RegisterDriver("mysql", orm.DRMySQL)
	//orm.RegisterDataBase("default", "mysql", "root:123456@tcp(localhost:3306)/lflss?charset=utf8")
	//orm.SetMaxIdleConns("default", 1000)
	//orm.SetMaxOpenConns("default", 2000)
	//
	//orm.RegisterModel(new(Teacher))
	db, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/lflss?charset=utf8")
	if err != nil {
		panic("cannot connect to mysql")
	}
	// init global
	Tm.Init()
	Cm.Init()
	Vm.Init()
}

type Teacher struct {
	ID     int
	Gender int    // 1: Male, 2: Female, 3: unknown
	Name   string // 姓名
	Mobile string // 手机号
}

type TeacherManager struct {
	nameMap map[string]*Teacher // 名字索引
	idMap   map[int]*Teacher    // id索引表
	nextID  int
	idMutex sync.Mutex
}

var (
	ErrInvalidParam = errors.New("invalid parameter")
	ErrNameExist    = errors.New("name exist")
	ErrNotExist     = errors.New("teacher not exist")
)

func (tm *TeacherManager) Init() {
	tm.idMap = make(map[int]*Teacher)
	tm.nameMap = make(map[string]*Teacher)
}

func (tm *TeacherManager) getID() int {
	tm.idMutex.Lock()
	tm.nextID++
	ret := tm.nextID
	tm.idMutex.Unlock()
	return ret
}

func (tm *TeacherManager) AddTeacher(t *Teacher) error {
	if t.ID != 0 {
		return ErrInvalidParam
	}

	if t.Name == "" {
		return ErrInvalidParam
	}

	if _, ok := tm.nameMap[t.Name]; ok {
		return ErrNameExist
	}

	if t.Mobile == "" {
		return ErrInvalidParam
	}

	if t.Gender <= 0 && t.Gender > 3 {
		return ErrInvalidParam
	}

	// assign new id
	t.ID = tm.getID()

	// add to map
	tm.nameMap[t.Name] = t
	tm.idMap[t.ID] = t
	fmt.Println("id=", t.ID)
	return nil
}

func (tm *TeacherManager) DelTeacher(id int) error {
	val, ok := tm.idMap[id]
	if !ok {
		return ErrNotExist
	}

	// remove from map
	delete(tm.nameMap, val.Name)
	delete(tm.idMap, val.ID)

	return nil
}

func (tm *TeacherManager) GetTeacherInfo(id int) (*Teacher, error) {
	ret := &Teacher{}
	err := ErrNotExist
	if val, ok := tm.idMap[id]; ok {
		*ret = *val
		err = nil
	}
	return ret, err
}

func (tm *TeacherManager) GetTeacherList(ids []int) ([]*Teacher, error) {
	ret := []*Teacher{}
	for _, v := range ids {
		if val, ok := tm.idMap[v]; ok {
			tmp := *val
			ret = append(ret, &tmp)
		}
	}
	return ret, nil
}

func (tm *TeacherManager) GetAll() []*Teacher {
	ret := []*Teacher{}
	for _, v := range tm.idMap {
		tmp := *v
		ret = append(ret, &tmp)
	}
	return ret
}
