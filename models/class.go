package models

import (
	"errors"
	"github.com/arong/dean/base"
	"sort"
	"sync"

	"github.com/astaxie/beego/logs"
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

func (cm *ClassManager) ModifyClass(c *Class) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	curr, ok := cm.idMap[c.ID]
	if !ok {
		return ErrClassNotExist
	}

	if curr.Equal(*c) {
		return nil
	}

	if c.Term != 0 {
		curr.Term = c.Term
	}

	if c.MasterID != 0 && c.MasterID != curr.MasterID {
		curr.MasterID = c.MasterID
	}

	if c.Year != 0 && c.Year != curr.Year {
		curr.Year = c.Year
	}

	err := Ma.UpdateClass(c)
	if err != nil {
		logs.Warn("database error")
		return err
	}

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

func (cm *ClassManager) GetAll() base.CommList {
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

func (cm *ClassManager) GetInfo(id int) (*ClassResp, error) {
	ret := &ClassResp{}

	if cm == nil {
		return ret, nil
	}

	val, ok := cm.idMap[id]
	if !ok {
		return ret, ErrClassNotExist
	}

	ret.Class = *val
	return ret, nil
}
