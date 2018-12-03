package models

import (
	"errors"
	"github.com/arong/dean/base"
	"github.com/astaxie/beego/logs"
	"sort"
	"sync"
)

// global manager
var Cm ClassManager

const (
	prefix = "高"
	suffix = "班"
)

var chineseNumberMap = map[int]string{
	1:  "一",
	2:  "二",
	3:  "三",
	4:  "四",
	5:  "五",
	6:  "六",
	7:  "七",
	8:  "八",
	9:  "九",
	10: "十",
}

var (
	// ErrClassNotExist class not exist
	ErrClassNotExist = errors.New("class not exist")
)

type ClassManager struct {
	idMap map[int]*Class
	mutex sync.Mutex
}

func (cm *ClassManager) Lock() {
	cm.mutex.Lock()
}

func (cm *ClassManager) UnLock() {
	cm.mutex.Unlock()
}

// maintain the relation between class and teacher
func (cm *ClassManager) Init(data map[int]*Class) {
	if cm == nil {
		return
	}

	if data == nil {
		cm.idMap = make(map[int]*Class)
	} else {
		cm.idMap = data
	}
}

func (cm *ClassManager) AddClass(c *Class) (int, error) {
	var ret int

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 查找主键重复
	if _, ok := cm.idMap[c.ID]; ok {
		logs.Debug("class id exist", c.ID)
		return ret, errExist
	}

	if c.Name == "" {
		c.Name = prefix + chineseNumberMap[c.Grade] + chineseNumberMap[c.Index] + suffix
	}

	err := Ma.InsertClass(c)
	if err != nil {
		logs.Warn("database error")
		return ret, err
	}

	cm.idMap[c.ID] = c

	logs.Info("create a new class", "classID", c.ID)
	return c.ID, nil
}

func (cm *ClassManager) ModifyClass(r *Class) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	curr, ok := cm.idMap[r.ID]
	if !ok {
		return ErrClassNotExist
	}

	err := Tm.CheckInstructorList(r.TeacherList)
	if err != nil {
		logs.Warn("[ModifyClass] CheckInstructorList error", "err", err)
		return err
	}

	if curr.Equal(*r) {
		logs.Debug("[ModifyClass] need do nothing")
		return nil
	}

	if r.Term != 0 {
		curr.Term = r.Term
	}

	if r.MasterID != 0 && r.MasterID != curr.MasterID {
		curr.MasterID = r.MasterID
	}

	if r.Year != 0 && r.Year != curr.Year {
		curr.Year = r.Year
	}

	// diff two list
	curr.TeacherList, curr.AddList, curr.RemoveList = curr.TeacherList.Diff(r.TeacherList)
	logs.Debug("[ModifyClass]", "addList", curr.AddList, "delList", curr.RemoveList, "all", curr.TeacherList)
	err = Ma.UpdateClass(curr)
	if err != nil {
		logs.Warn("[ModifyClass] database error")
		return err
	}
	curr.AddList = InstructorList{}
	curr.RemoveList = InstructorList{}
	return nil
}

func (cm *ClassManager) DelClass(list ClassIDList) (ClassIDList, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	failedList := ClassIDList{}

	for _, id := range list {
		_, ok := cm.idMap[id]
		if !ok {
			failedList = append(failedList, id)
			continue
		}

		err := Ma.DeleteClass(id)
		if err != nil {
			logs.Warn("[DelClass] database failed", "err", err)
			failedList = append(failedList, id)
		}
		delete(cm.idMap, id)
	}

	return failedList, nil
}

func (cm *ClassManager) Filter() base.CommList {
	resp := base.CommList{}
	list := ClassList{}
	for _, v := range cm.idMap {
		list = append(list, v)
	}
	sort.Sort(list)
	resp.List = list
	resp.Total = len(list)
	return resp
}

func (cm *ClassManager) GetAll() ItemList {
	resp := ItemList{}
	for _, v := range cm.idMap {
		resp = append(resp, Item{ID: v.ID, Name: v.Name})
	}
	sort.Sort(resp)
	return resp
}

func (cm *ClassManager) GetInfo(id int) (*Class, error) {
	ret := &Class{}

	if cm == nil {
		return ret, nil
	}

	val, ok := cm.idMap[id]
	if !ok {
		return ret, ErrClassNotExist
	}

	*ret = *val
	return ret, nil
}
