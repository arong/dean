package manager

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/arong/dean/models"

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

type TeacherManager struct {
	// list is the ultimate list to hold all teacher info
	// I choose list over map solely because list is good for
	// iterator over and filter, it cause less cache miss
	// add will always append item at tail
	// delete will just set the status to be deleted
	list models.TeacherList
	// name map hold the name to access teacher info
	nameMap map[string]int // 名字索引
	// id map hold id to access teacher info
	idMap map[int64]int // id索引表
	// protect inner data with mutex
	mutex sync.Mutex // 保护前两个数据表
	// save deleted item count
	deletedCount int32
	// signal channel
	ch    chan bool
	store models.TeacherStore
}

func (tm *TeacherManager) save(t models.Teacher) {
	t.Status = base.StatusValid
	tm.list = append(tm.list, t)
	k := len(tm.list) - 1
	tm.nameMap[t.Name] = k
	tm.idMap[t.TeacherID] = k
}

func (tm *TeacherManager) get(id int64) (models.Teacher, error) {
	k, ok := tm.idMap[id]
	if !ok {
		return models.Teacher{}, errNotExist
	}
	return tm.list[k], nil
}

func (tm *TeacherManager) update(t models.Teacher) {
	k, ok := tm.idMap[t.TeacherID]
	if !ok {
		logs.Warn("bug found")
		return
	}
	tm.list[k] = t
}

func (tm *TeacherManager) delete(k int) {
	p := &tm.list[k]
	p.Status = base.StatusDeleted
	delete(tm.nameMap, p.Name)
	delete(tm.idMap, p.TeacherID)
	curr := atomic.AddInt32(&tm.deletedCount, 1)

	if len(tm.list) > 0 && float64(curr)/float64(len(tm.list)) >= 0.75 {
		tm.ch <- true
	}
}

func (tm *TeacherManager) doClean() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if len(tm.list) > 0 && float64(tm.deletedCount)/float64(len(tm.list)) >= 0.75 {
		list := models.TeacherList{}
		idMap := make(map[int64]int)
		nameMap := make(map[string]int)

		for _, v := range tm.list {
			if v.Status == base.StatusDeleted {
				continue
			}

			list = append(list, v)
			k := len(list) - 1
			idMap[v.TeacherID] = k
			nameMap[v.Name] = k
		}

		tm.list = list
		tm.idMap = idMap
		tm.nameMap = nameMap
		// reset delete count
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
func (tm *TeacherManager) Init(data models.TeacherList) {
	tm.nameMap = make(map[string]int)
	tm.idMap = make(map[int64]int)
	tm.list = models.TeacherList{}

	if tm == nil || data == nil {
		return
	}
	sort.Sort(data)

	tm.mutex.Lock()

	for k, v := range data {
		if v.Status != base.StatusValid {
			continue
		}
		tmp := v
		tm.list = append(tm.list, tmp)
		tm.nameMap[v.Name] = k
		tm.idMap[v.TeacherID] = k
		if v.SubjectID > 0 {
			// todo: fixup the tie
			//Sm.IncRef(v.SubjectID)
		}
	}

	tm.mutex.Unlock()
	tm.ch = make(chan bool)
	go tm.clean()
}

// AddTeacher: AddTeacher
func (tm *TeacherManager) AddTeacher(t *models.Teacher) (int64, error) {
	// suppress the decoded value
	t.TeacherID = 0

	if _, ok := tm.nameMap[t.Name]; ok {
		return 0, errNameExist
	}

	err := t.Check()
	if err != nil {
		logs.Warn("[TeacherManager::AddTeacher] invalid parameter", "err", err)
		return 0, err
	}

	t.TeacherID, err = tm.store.SaveTeacher(*t)
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

	// todo: fixup this
	//if t.SubjectID > 0 {
	//	Sm.IncRef(t.SubjectID)
	//}
	return t.TeacherID, nil
}

func (tm *TeacherManager) UpdateTeacher(t *models.Teacher) error {
	if t.TeacherID <= 0 {
		logs.Info("[TeacherManager::UpdateTeacher] invalid parameter")
		return errInvalidParam
	}

	err := t.Check()
	if err != nil {
		logs.Info("[TeacherManager::UpdateTeacher] invalid parameter", "err", err)
		return err
	}

	curr, err := tm.get(t.TeacherID)
	if err != nil {
		logs.Info("[TeacherManager::UpdateTeacher] id not found", "err", err)
		return err
	}

	if curr.Equal(t.TeacherMeta) {
		logs.Info("[TeacherManager::UpdateTeacher] nothing to do")
		return nil
	}

	// accept change
	if curr.Name != t.Name {
		return errors.New("[TeacherManager::UpdateTeacher] name could not change")
	}

	// subject
	oldSubjectID := curr.SubjectID
	if curr.SubjectID != t.SubjectID {
		curr.SubjectID = t.SubjectID
	}

	if curr.Mobile != t.Mobile {
		curr.Mobile = t.Mobile
	}

	if curr.Birthday != t.Birthday {
		curr.Birthday = t.Birthday
		curr.Age = t.Age
	}

	if curr.Address != t.Address {
		curr.Address = t.Address
	}

	err = tm.store.UpdateTeacher(curr)
	if err != nil {
		logs.Info("[TeacherManager::UpdateTeacher] UpdateTeacher failed", "err", err)
		return err
	}

	{
		tm.mutex.Lock()
		tm.update(curr)
		tm.mutex.Unlock()
	}

	if oldSubjectID != curr.SubjectID {
		Sm.IncRef(curr.SubjectID)
		Sm.DecRef(oldSubjectID)
	}

	return nil
}

// DelTeacher: DelTeacher
func (tm *TeacherManager) DelTeacher(idList []int64) ([]int64, error) {
	tm.mutex.Lock()
	tm.mutex.Unlock()

	failed := []int64{}
	list := []int64{}
	for _, v := range idList {
		k, ok := tm.idMap[v]
		if !ok {
			logs.Debug("[TeacherManager::DelTeacher] not found", "id", v)
			failed = append(failed, v)
			continue
		}

		list = append(list, v)
		// remove from map
		tm.delete(k)
	}

	// delete from database
	if len(list) > 0 {
		err := tm.store.DeleteTeacher(idList)
		if err != nil {
			failed = append(failed, list...)
			logs.Warn("[TeacherManager::DelTeacher] database error")
		}
	}
	return failed, nil
}

type TeacherInfoResp struct {
	models.TeacherMeta
	SubjectID int `json:"subject_id"`
}

// GetTeacherInfo: GetTeacherInfo
func (tm *TeacherManager) GetTeacherInfo(id int64) (TeacherInfoResp, error) {
	ret := TeacherInfoResp{}
	err := errNotExist
	tm.mutex.Lock()
	if val, ok := tm.idMap[id]; ok {
		p := tm.list[val]
		ret.TeacherMeta = p.TeacherMeta
		err = nil
	}
	tm.mutex.Unlock()

	if ret.SubjectID > 0 {
		ret.Subject = Sm.getSubjectName(ret.SubjectID)
	}
	return ret, err
}

//type TeacherFilter struct {
//	base.CommPage
//	Gender int    `json:"gender"`
//	Age    int    `json:"age"`
//	Name   string `json:"name"`
//	Mobile string `json:"mobile"`
//}

func (tm *TeacherManager) Filter(f models.TeacherFilter) base.CommList {
	ret := models.TeacherList{}

	tm.mutex.Lock()
	total, ret := tm.list.Filter(f)
	tm.mutex.Unlock()

	logs.Debug("[Filter]", "total", total)

	return base.CommList{Total: total, List: ret}
}

type simpleTeacher struct {
	ID   int64  `json:"id"`
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

func (tm *TeacherManager) GetAll() base.CommList {
	ret := simpleTeacherList{}
	tm.mutex.Lock()
	for _, v := range tm.list {
		if v.Status == base.StatusDeleted {
			continue
		}
		ret = append(ret, simpleTeacher{Name: v.Name, ID: v.TeacherID})
	}
	tm.mutex.Unlock()
	sort.Sort(ret)
	return base.CommList{Total: len(ret), List: ret}
}

func (tm *TeacherManager) IsTeacherExist(id int64) bool {
	tm.mutex.Lock()
	_, ok := tm.idMap[id]
	tm.mutex.Unlock()
	return ok
}

//func (tm *TeacherManager) CheckTeachers(ids []int64) bool {
//	tm.mutex.Lock()
//	defer tm.mutex.Unlock()
//
//	for _, v := range ids {
//		if _, ok := tm.idMap[v]; !ok {
//			return false
//		}
//	}
//	return true
//}

// CheckInstructorList check to see if the id in list exist
func (tm *TeacherManager) CheckInstructorList(list models.InstructorList) error {
	for _, v := range list {
		if t, ok := tm.idMap[v.TeacherID]; ok {
			p := &tm.list[t]
			if p.SubjectID != v.SubjectID {
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
