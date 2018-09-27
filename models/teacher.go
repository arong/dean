package models

import (
	"errors"
	"fmt"
	"sync"
)

// a global handler
var Tm TeacherManager

//var db *sql.DB

func init() {
	// init global
	Ma.Init()
	Tm.Init()
	Cm.Init()
	Vm.Init()
}

type Teacher struct {
	ID     int64
	Gender int    // 1: Male, 2: Female, 3: unknown
	Name   string // 姓名
	Mobile string // 手机号
}

type TeacherList []*Teacher

func (tl TeacherList) Len() int {
	return len(tl)
}

func (tl TeacherList) Swap(i, j int) {
	tl[i], tl[j] = tl[i], tl[j]
}

func (tl TeacherList) Less(i, j int) bool {
	return tl[i].ID < tl[j].ID
}

type TeacherManager struct {
	nameMap map[string]*Teacher // 名字索引
	idMap   map[int64]*Teacher  // id索引表
	nextID  int64
	idMutex sync.Mutex
}

var (
	ErrInvalidParam = errors.New("invalid parameter")
	ErrNameExist    = errors.New("name exist")
	ErrNotExist     = errors.New("teacher not exist")
)

func (tm *TeacherManager) Init() {
	tm.idMap = make(map[int64]*Teacher)
	tm.nameMap = make(map[string]*Teacher)

	all, err := Ma.loadAllTeachers()
	if err != nil {
		return
	}

	for _, v := range all {
		tm.idMap[v.ID] = v
		tm.nameMap[v.Name] = v
		if tm.nextID < v.ID {
			tm.nextID = v.ID
		}
	}
}

func (tm *TeacherManager) getID() int64 {
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
	//t.ID = tm.getID()

	err := Ma.InsertTeacher(t)
	if err != nil {
		return err
	}
	// add to map
	tm.nameMap[t.Name] = t
	tm.idMap[t.ID] = t
	fmt.Println("id=", t.ID)
	return nil
}

func (tm *TeacherManager) DelTeacher(id int64) error {
	val, ok := tm.idMap[id]
	if !ok {
		return ErrNotExist
	}

	// remove from map
	delete(tm.nameMap, val.Name)
	delete(tm.idMap, val.ID)

	return nil
}

func (tm *TeacherManager) GetTeacherInfo(id int64) (*Teacher, error) {
	ret := &Teacher{}
	err := ErrNotExist
	if val, ok := tm.idMap[id]; ok {
		*ret = *val
		err = nil
	}
	return ret, err
}

func (tm *TeacherManager) GetTeacherList(ids []int64) ([]*Teacher, error) {
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
