package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"os"
	"sort"
	"sync"
)

const (
	eGenderMale    = 1
	eGenderFemale  = 2
	eGenderUnknown = 3
)

// a global handler
var Tm TeacherManager

func Init(conf *DBConfig) {
	// allocate memory
	Ma.Init(conf)
	Vm.Init()

	// data warm up
	err := Ma.LoadAllData()
	if err != nil {
		logs.Error("init failed", "err", err)
		os.Exit(-1)
	}
}

type Teacher struct {
	SubjectID int
	profile
}

type TeacherList []*Teacher

func (tl TeacherList) Len() int {
	return len(tl)
}

func (tl TeacherList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (tl TeacherList) Less(i, j int) bool {
	return tl[i].ID < tl[j].ID
}

type TeacherManager struct {
	nameMap map[string]*Teacher // 名字索引
	idMap   map[UserID]*Teacher // id索引表
	mutex   sync.Mutex          // 保护前两个数据表
}

var (
	ErrInvalidParam = errors.New("invalid parameter")
	ErrNameExist    = errors.New("name exist")
	ErrNotExist     = errors.New("teacher not exist")
)

func (tm *TeacherManager) Init(data map[UserID]*Teacher) {
	if tm == nil {
		return
	}

	tm.nameMap = make(map[string]*Teacher)

	if data == nil {
		tm.idMap = make(map[UserID]*Teacher)
	} else {
		tm.idMap = data
		for _, v := range data {
			tm.nameMap[v.Name] = v
		}
	}
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

	if t.Gender < eGenderMale && t.Gender > eGenderUnknown {
		return ErrInvalidParam
	}

	err := Ma.InsertTeacher(t)
	if err != nil {
		logs.Warn("[] database error")
		return err
	}

	// add to map
	{
		tm.mutex.Lock()
		tm.nameMap[t.Name] = t
		tm.idMap[t.ID] = t
		tm.mutex.Unlock()
	}

	logs.Info("id=", t.ID)
	return nil
}

func (tm *TeacherManager) DelTeacher(id UserID) error {
	tm.mutex.Lock()
	tm.mutex.Unlock()

	val, ok := tm.idMap[id]
	if !ok {
		logs.Debug("[TeacherManager::DelTeacher] not found", "id", id)
		return ErrNotExist
	}

	// delete from database
	err := Ma.DeleteTeacher(id)
	if err != nil {
		logs.Warn("[TeacherManager::DelTeacher] database error")
		return err
	}

	// remove from map
	delete(tm.nameMap, val.Name)
	delete(tm.idMap, val.ID)

	return nil
}

func (tm *TeacherManager) GetTeacherInfo(id UserID) (*Teacher, error) {
	ret := &Teacher{}
	err := ErrNotExist
	tm.mutex.Lock()
	if val, ok := tm.idMap[id]; ok {
		*ret = *val
		err = nil
	}
	tm.mutex.Unlock()
	return ret, err
}

func (tm *TeacherManager) GetTeacherList(ids []UserID) (TeacherList, error) {
	ret := TeacherList{}
	tm.mutex.Lock()
	for _, v := range ids {
		if val, ok := tm.idMap[v]; ok {
			tmp := *val
			ret = append(ret, &tmp)
		}
	}
	tm.mutex.Unlock()
	//sort.Sort(ret)
	return ret, nil
}

func (tm *TeacherManager) GetAll() TeacherList {
	ret := TeacherList{}
	tm.mutex.Lock()
	for _, v := range tm.idMap {
		tmp := *v
		ret = append(ret, &tmp)
	}
	tm.mutex.Unlock()
	sort.Sort(ret)
	return ret
}

func (tm *TeacherManager) IsTeacherExist(id UserID) bool {
	tm.mutex.Lock()
	_, ok := tm.idMap[id]
	tm.mutex.Unlock()
	return ok
}

func (tm *TeacherManager) CheckTeachers(ids []UserID) bool {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	for _, v := range ids {
		if _, ok := tm.idMap[v]; !ok {
			return false
		}
	}
	return true
}
