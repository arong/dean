package models

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

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
	// store is the ultimate list to hold all teacher info
	// I choose list over map solely because list is good for
	// iterator over and filter, it cause less cache miss
	// add will always append item at tail
	// delete will just set the status to be deleted
	store TeacherList
	// name map hold the name to access teacher info
	nameMap map[string]int // 名字索引
	// id map hold id to access teacher info
	idMap map[int64]int // id索引表
	// protect inner data with mutex
	mutex sync.Mutex // 保护前两个数据表
	// save deleted item count
	deletedCount int32
	// signal channel
	ch chan bool
}

func (tm *TeacherManager) save(t Teacher) {
	t.Status = base.StatusValid
	tm.store = append(tm.store, t)
	k := len(tm.store) - 1
	tm.nameMap[t.Name] = k
	tm.idMap[t.TeacherID] = k
}

func (tm *TeacherManager) update(t Teacher) {
	k, ok := tm.idMap[t.TeacherID]
	if !ok {
		return
	}
	tm.store[k] = t
}

func (tm *TeacherManager) delete(k int) {
	p := &tm.store[k]
	p.Status = base.StatusDeleted
	delete(tm.nameMap, p.Name)
	delete(tm.idMap, p.TeacherID)
	curr := atomic.AddInt32(&tm.deletedCount, 1)

	if len(tm.store) > 0 && float64(curr)/float64(len(tm.store)) >= 0.75 {
		logs.Warn("fire a delete signal")
		tm.ch <- true
	}
}

func (tm *TeacherManager) doClean() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if float64(tm.deletedCount)/float64(len(tm.store)) >= 0.75 {
		list := TeacherList{}
		for _, v := range tm.store {
			if v.Status == base.StatusDeleted {
				continue
			}
			list = append(list, v)
		}
		idMap := make(map[int64]int)
		nameMap := make(map[string]int)
		for k, v := range list {
			idMap[v.TeacherID] = k
			nameMap[v.Name] = k
		}

		tm.store = list
		tm.idMap = idMap
		tm.nameMap = nameMap
		atomic.StoreInt32(&tm.deletedCount, 0)
	}
}

func (tm *TeacherManager) clean() {
	t := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-t.C:
			logs.Warn("scheduled clean")
			tm.doClean()
		case <-tm.ch:
			logs.Warn("triggered clean")
			tm.doClean()
		}
	}
}

var (
	errInvalidParam = errors.New("invalid parameter")
	errNameExist    = errors.New("name exist")
)

// Init: Init
func (tm *TeacherManager) Init(data TeacherList) {
	if tm == nil || data == nil {
		return
	}
	tm.mutex.Lock()

	sort.Sort(data)
	tm.store = data

	tm.nameMap = make(map[string]int)
	tm.idMap = make(map[int64]int)
	for k, v := range data {
		tm.nameMap[v.Name] = k
		tm.idMap[v.TeacherID] = k
	}

	tm.mutex.Unlock()
	tm.ch = make(chan bool)
	go tm.clean()
}

// AddTeacher: AddTeacher
func (tm *TeacherManager) AddTeacher(t *Teacher) (int64, error) {
	if t.TeacherID != 0 {
		return 0, errInvalidParam
	}

	if _, ok := tm.nameMap[t.Name]; ok {
		return 0, errNameExist
	}

	err := t.Check()
	if err != nil {
		logs.Warn("[TeacherManager::AddTeacher] invalid parameter", "err", err)
		return 0, err
	}

	t.TeacherID, err = Ma.InsertTeacher(*t)
	if err != nil {
		logs.Warn("[TeacherManager::AddTeacher] database error")
		return 0, err
	}

	// add to map
	{
		tm.mutex.Lock()
		tm.save(*t)
		tm.mutex.Unlock()
	}

	return t.TeacherID, nil
}

// ModTeacher: ModTeacher
func (tm *TeacherManager) ModTeacher(t *Teacher) error {
	if t.TeacherID == 0 {
		logs.Info("[TeacherManager::ModTeacher] invalid parameter")
		return errInvalidParam
	}

	err := t.Check()
	if err != nil {
		logs.Info("[TeacherManager::ModTeacher] invalid parameter", "err", err)
		return err
	}

	err = Ma.UpdateTeacher(*t)
	if err != nil {
		logs.Info("[TeacherManager::ModTeacher] UpdateTeacher failed", "err", err)
		return err
	}

	{
		tm.mutex.Lock()
		tm.update(*t)
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
		k, ok := tm.idMap[v]
		if !ok {
			logs.Debug("[TeacherManager::DelTeacher] not found", "id", v)
			failed = append(failed, v)
			continue
		}

		// remove from map
		tm.delete(k)
	}
	// delete from database
	err := Ma.DeleteTeacher(idList)
	if err != nil {
		logs.Warn("[TeacherManager::DelTeacher] database error")
	}
	return failed, nil
}

// GetTeacherInfo: GetTeacherInfo
func (tm *TeacherManager) GetTeacherInfo(id int64) (*TeacherInfoResp, error) {
	ret := &TeacherInfoResp{}
	err := errNotExist
	tm.mutex.Lock()
	if val, ok := tm.idMap[id]; ok {
		p := tm.store[val]
		ret.TeacherMeta = p.TeacherMeta
		err = nil
	}
	tm.mutex.Unlock()

	ret.Subject = Sm.getSubjectName(ret.SubjectID)
	return ret, err
}

// Filter: Filter
func (tm *TeacherManager) Filter(f *TeacherFilter) base.CommList {
	ret := TeacherList{}

	tm.mutex.Lock()
	list := IntList{}
	for k, v := range tm.store {
		if v.Status != base.StatusValid {
			continue
		}

		if f.Gender != 0 && f.Gender != v.Gender {
			continue
		}

		if f.Name != "" && f.Name != v.Name {
			continue
		}
		list = append(list, k)
	}

	total := len(list)
	list = list.Page(f.CommPage)
	for _, v := range list {
		ret = append(ret, tm.store[v])
	}

	logs.Debug("[TeacherManager::GetAll]", "len(store)", len(tm.store))
	tm.mutex.Unlock()

	return base.CommList{Total: total, List: ret}
}

// GetAll: GetAll
func (tm *TeacherManager) GetAll() base.CommList {
	ret := simpleTeacherList{}
	tm.mutex.Lock()
	for _, v := range tm.store {
		if v.Status == base.StatusDeleted {
			continue
		}
		ret = append(ret, simpleTeacher{Name: v.Name, ID: v.TeacherID})
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
			p := &tm.store[t]
			if p.SubjectID != 0 && p.SubjectID != v.SubjectID {
				logs.Debug("[CheckInstructorList]", "t.SubjectID", p.SubjectID, "v.SubjectID", v.SubjectID)
				return fmt.Errorf("teacher %s and subject %d not match", p.Name, v.SubjectID)
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
