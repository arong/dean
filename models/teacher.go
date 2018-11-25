package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"os"
	"sort"
	"sync"
	"time"
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
	TeacherID UserID `json:"teacher_id"`
	SubjectID int    `json:"subject_id"`
	profile
}

func (t *Teacher) IsValid() error {
	if t.Mobile == "" {
		return errors.New("invalid mobile")
	}

	if t.Gender < eGenderMale && t.Gender > eGenderUnknown {
		return errors.New("invalid gender")
	}

	if t.Birthday == "" {
		return errors.New("empty birthday")
	}

	_, err := time.Parse("2006-01-02", t.Birthday)
	if err != nil {
		return errors.New("invalid birthday")
	}
	return nil
}

type TeacherList []*Teacher

func (tl TeacherList) Len() int {
	return len(tl)
}

func (tl TeacherList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (tl TeacherList) Less(i, j int) bool {
	return tl[i].TeacherID < tl[j].TeacherID
}

type TeacherListResp struct {
	Total      int         `json:"total"`
	RecordList interface{} `json:"list"`
}

type TeacherFilter struct {
	Page   int    `json:"page"` // start from 1
	Size   int    `json:"size"`
	Gender int    `json:"gender"`
	Age    int    `json:"age"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
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

func (tm *TeacherManager) Lock() {
	tm.mutex.Lock()
}

func (tm *TeacherManager) UnLock() {
	tm.mutex.Unlock()
}

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
			tm.nameMap[v.RealName] = v
		}
	}
}

func (tm *TeacherManager) AddTeacher(t *Teacher) error {
	if t.TeacherID != 0 {
		return ErrInvalidParam
	}

	if _, ok := tm.nameMap[t.RealName]; ok {
		return ErrNameExist
	}

	err := t.IsValid()
	if err != nil {
		logs.Warn("[TeacherManager::AddTeacher] invalid parameter", "err", err)
		return err
	}

	err = Ma.InsertTeacher(t)
	if err != nil {
		logs.Warn("[TeacherManager::AddTeacher] database error")
		return err
	}

	// add to map
	{
		tm.mutex.Lock()
		tm.nameMap[t.RealName] = t
		tm.idMap[t.TeacherID] = t
		tm.mutex.Unlock()
	}

	logs.Info("id=", t.TeacherID)
	return nil
}

func (tm *TeacherManager) ModTeacher(t *Teacher) error {
	if t.TeacherID == 0 {
		logs.Info("[TeacherManager::ModTeacher] invalid parameter")
		return ErrInvalidParam
	}

	err := t.IsValid()
	if err != nil {
		logs.Info("[TeacherManager::ModTeacher] invalid parameter", "err", err)
		return err
	}

	err = Ma.UpdateTeacher(t)
	if err != nil {
		logs.Info("[TeacherManager::ModTeacher] UpdateTeacher failed", "err", err)
		return err
	}

	{
		tm.mutex.Lock()
		tm.nameMap[t.RealName] = t
		tm.idMap[t.TeacherID] = t
		tm.mutex.Unlock()
	}
	return nil
}

func (tm *TeacherManager) DelTeacher(idList []UserID) ([]UserID, error) {
	tm.mutex.Lock()
	tm.mutex.Unlock()

	failed := []UserID{}
	for _, v := range idList {
		val, ok := tm.idMap[v]
		if !ok {
			logs.Debug("[TeacherManager::DelTeacher] not found", "id", v)
			failed = append(failed, v)
			continue
		}

		// delete from database
		err := Ma.DeleteTeacher(v)
		if err != nil {
			logs.Warn("[TeacherManager::DelTeacher] database error")
			failed = append(failed, v)
			continue
		}

		// remove from map
		delete(tm.nameMap, val.RealName)
		delete(tm.idMap, val.TeacherID)
	}
	return failed, nil
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

func (tm *TeacherManager) Filter(f *TeacherFilter) TeacherListResp {
	ret := TeacherList{}
	total := 0
	start := (f.Page - 1) * f.Size
	end := f.Page * f.Size

	logs.Debug("[TeacherManager::GetAll]", "start", start, "end", end)
	tm.mutex.Lock()
	for _, v := range tm.idMap {
		if f.Gender != 0 && f.Gender != v.Gender {
			continue
		}

		if f.Name != "" && f.Name != v.RealName {
			continue
		}

		total++
		tmp := v
		ret = append(ret, tmp)
	}

	sort.Sort(ret)
	sub := TeacherList{}
	for k, v := range ret {
		if k < start || k >= end {
			continue
		}
		tmp := v
		sub = append(sub, tmp)
	}
	tm.mutex.Unlock()
	return TeacherListResp{Total: total, RecordList: sub}
}

type simpleTeacher struct {
	ID   UserID `json:"id"`
	Name string `json:"name"`
}

type simpleTeacherList []simpleTeacher

func (s simpleTeacherList) Len() int {
	return len(s)
}
func (s simpleTeacherList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s simpleTeacherList) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}
func (tm *TeacherManager) GetAll() TeacherListResp {
	ret := simpleTeacherList{}
	tm.mutex.Lock()
	for _, v := range tm.idMap {
		ret = append(ret, simpleTeacher{Name: v.RealName, ID: v.TeacherID})
	}
	tm.mutex.Unlock()
	sort.Sort(ret)
	return TeacherListResp{Total: len(ret), RecordList: ret}
}

func (tm *TeacherManager) FilterTeacher() (TeacherListResp, error) {
	ret := TeacherListResp{}

	return ret, nil
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

func (tm *TeacherManager) CheckInstructorList(list InstructorList) bool {
	for _, v := range list {
		if t, ok := tm.idMap[v.TeacherID]; ok {
			if t.SubjectID != v.SubjectID {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (tm *TeacherManager) CheckID(id UserID) bool {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	_, ok := tm.idMap[id]
	return ok
}
