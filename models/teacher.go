package models

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/arong/dean/base"

	"github.com/astaxie/beego/logs"
)

const (
	eGenderMale    = 1
	eGenderFemale  = 2
	eGenderUnknown = 3
)

// Tm a global handler
var Tm TeacherManager

// Init init da config
func Init(conf *DBConfig) {
	// allocate memory
	Ma.Init(conf)

	// data warm up
	err := Ma.LoadAllData()
	if err != nil {
		logs.Error("init failed", "err", err)
		os.Exit(-1)
	}
}

type TeacherManager struct {
	nameMap map[string]*Teacher // 名字索引
	idMap   map[int64]*Teacher  // id索引表
	mutex   sync.Mutex          // 保护前两个数据表
}

var (
	errInvalidParam = errors.New("invalid parameter")
	errNameExist    = errors.New("name exist")
)

// Init: Init
func (tm *TeacherManager) Init(data map[int64]*Teacher) {
	if tm == nil {
		return
	}

	tm.nameMap = make(map[string]*Teacher)

	if data == nil {
		tm.idMap = make(map[int64]*Teacher)
	} else {
		tm.idMap = data
		for _, v := range data {
			tm.nameMap[v.RealName] = v
		}
	}
}

// AddTeacher: AddTeacher
func (tm *TeacherManager) AddTeacher(t *Teacher) error {
	if t.TeacherID != 0 {
		return errInvalidParam
	}

	if _, ok := tm.nameMap[t.RealName]; ok {
		return errNameExist
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

// ModTeacher: ModTeacher
func (tm *TeacherManager) ModTeacher(t *Teacher) error {
	if t.TeacherID == 0 {
		logs.Info("[TeacherManager::ModTeacher] invalid parameter")
		return errInvalidParam
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

// DelTeacher: DelTeacher
func (tm *TeacherManager) DelTeacher(idList []int64) ([]int64, error) {
	tm.mutex.Lock()
	tm.mutex.Unlock()

	failed := []int64{}
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

// GetTeacherInfo: GetTeacherInfo
func (tm *TeacherManager) GetTeacherInfo(id int64) (*Teacher, error) {
	ret := &Teacher{}
	err := errNotExist
	tm.mutex.Lock()
	if val, ok := tm.idMap[id]; ok {
		*ret = *val
		err = nil
	}
	tm.mutex.Unlock()

	ret.Subject = Sm.getSubjectName(ret.SubjectID)
	return ret, err
}

// GetTeacherList: GetTeacherList
func (tm *TeacherManager) GetTeacherList(ids []int64) (TeacherList, error) {
	ret := TeacherList{}
	tm.mutex.Lock()
	for _, v := range ids {
		if val, ok := tm.idMap[v]; ok {
			tmp := *val
			ret = append(ret, &tmp)
		}
	}
	tm.mutex.Unlock()
	return ret, nil
}

// Filter: Filter
func (tm *TeacherManager) Filter(f *TeacherFilter) base.CommList {
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

	tm.mutex.Unlock()
	sort.Sort(ret)
	ret = ret.Page(f.Page, f.Size)
	for k, v := range ret {
		ret[k].Subject = Sm.getSubjectName(v.SubjectID)
	}
	return base.CommList{Total: total, List: ret}
}

// GetAll: GetAll
func (tm *TeacherManager) GetAll() base.CommList {
	ret := simpleTeacherList{}
	tm.mutex.Lock()
	for _, v := range tm.idMap {
		ret = append(ret, simpleTeacher{Name: v.RealName, ID: v.TeacherID})
	}
	tm.mutex.Unlock()
	sort.Sort(ret)
	return base.CommList{Total: len(ret), List: ret}
}

// IsTeacherExist IsTeacherExist
func (tm *TeacherManager) IsTeacherExist(id int64) bool {
	tm.mutex.Lock()
	_, ok := tm.idMap[id]
	tm.mutex.Unlock()
	return ok
}

// CheckTeachers: CheckTeachers
func (tm *TeacherManager) CheckTeachers(ids []int64) bool {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	for _, v := range ids {
		if _, ok := tm.idMap[v]; !ok {
			return false
		}
	}
	return true
}

// CheckInstructorList check to see if the id in list exist
func (tm *TeacherManager) CheckInstructorList(list InstructorList) error {
	for _, v := range list {
		if t, ok := tm.idMap[v.TeacherID]; ok {
			if t.SubjectID != 0 && t.SubjectID != v.SubjectID {
				logs.Debug("[CheckInstructorList]", "t.SubjectID", t.SubjectID, "v.SubjectID", v.SubjectID)
				return fmt.Errorf("teacher %s and subject %d not match", t.RealName, v.SubjectID)
			}
		} else {
			return errors.New(fmt.Sprintf("teacher %d not found", v.TeacherID))
		}
	}
	return nil
}

// IsExist check to see if the teacher id exist
func (tm *TeacherManager) IsExist(id int64) bool {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	_, ok := tm.idMap[id]
	return ok
}
