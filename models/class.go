package models

import (
	"errors"
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
	idMap map[ClassID]*Class
	mutex sync.Mutex
}

// maintain the relation between class and teacher
func (cm *ClassManager) Init(data map[ClassID]*Class) {
	if cm == nil {
		return
	}

	if data == nil {
		cm.idMap = make(map[ClassID]*Class)
	} else {
		cm.idMap = data
	}
}

func (cm *ClassManager) AddClass(c *Class) (ClassID, error) {
	var ret ClassID

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

	if curr.Name != c.Name {
		curr.Name = c.Name

		err := Ma.UpdateClass(c)
		if err != nil {
			logs.Warn("database error")
			return err
		}
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

func (cm *ClassManager) GetAll() ClassList {
	ret := ClassList{}
	for _, v := range cm.idMap {
		ret = append(ret, v)
	}
	return ret
}

func (cm *ClassManager) GetInfo(id ClassID) (*ClassResp, error) {
	ret := &ClassResp{}
	var err error

	if cm == nil {
		return ret, nil
	}

	val, ok := cm.idMap[id]
	if !ok {
		return ret, ErrClassNotExist
	}

	ret.Class = *val
	//ret.Teachers, err = Tm.GetTeacherList(ret.TeacherIDs)
	if err != nil {
		return ret, err
	}
	return ret, nil
}
