package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"strconv"
	"sync"
)

// global manager
var Cm ClassManager

const (
	prefix = "高"
	suffix = "班"
)

var (
	ErrClassNotExist = errors.New("class not exist")
)

// of size two byte
// +--------+-------+
// | Grade  | index |
// +--------+-------+
type ClassID int

// Class is the
type Class struct {
	Filter
	ID         ClassID `json:"id"`
	Name       string  `json:"name"`
	TeacherIDs []UserID `json:"-"` // id
}

type ClassResp struct {
	Class
	Teachers TeacherList // 详情
}
type ClassList []*Class

func (cl ClassList) Len() int {
	return len(cl)
}

func (cl ClassList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

func (cl ClassList) Less(i, j int) bool {
	if cl[i].Grade < cl[j].Grade {
		return true
	} else if cl[i].Grade > cl[j].Grade {
		return false
	} else {
		return cl[i].Index < cl[j].Index
	}
}

type Filter struct {
	Grade int `json:"grade"` // 年级
	Index int `json:"index"` // 班级
}

func (f *Filter) GetID() ClassID {
	if f == nil {
		return 0
	}
	return ClassID(((f.Grade & 0xf) << 8) | f.Index&0x0f)
}

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
	c.ID = c.GetID()

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 查找主键重复
	if _, ok := cm.idMap[c.ID]; ok {
		logs.Debug("class id exist", c.ID)
		return ret, errExist
	}

	if c.Name == "" {
		c.Name = prefix + strconv.Itoa(c.Grade) + "-" + strconv.Itoa(c.Index) + suffix
	}

	if !Tm.CheckTeachers(c.TeacherIDs) {
		return ret, errNotExist
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

	// 检查教师是否有增减
	currMap := make(map[UserID]bool)
	for _, v := range curr.TeacherIDs {
		currMap[v] = true
	}

	newMap := make(map[UserID]bool)
	for _, v := range c.TeacherIDs {
		if _, ok := currMap[v]; ok {
			delete(currMap, v)
			continue
		}
		newMap[v] = true
	}

	if len(currMap) > 0 {
		addList := make([]UserID, len(newMap))
		for k, _ := range newMap {
			addList = append(addList, k)
		}

		err := Ma.AttachTeacher(c.ID, addList)
		if err != nil {
			logs.Warn("database error")
			return err
		}
	}

	if len(currMap) > 0 {
		removeList := make([]UserID, len(currMap))
		for k, _ := range currMap {
			removeList = append(removeList, k)
		}

		err := Ma.DetachTeacher(c.ID, removeList)
		if err != nil {
			logs.Warn("database error")
			return err
		}
	}
	return nil
}

func (cm *ClassManager) DelClass(id ClassID) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	_, ok := cm.idMap[id]
	if !ok {
		return ErrClassNotExist
	}

	err := Ma.DeleteClass(id)
	if err != nil {
		logs.Warn("[] database failed", "err", err)
		return err
	}

	delete(cm.idMap, id)

	return nil
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
	ret.Teachers, err = Tm.GetTeacherList(ret.TeacherIDs)
	if err != nil {
		return ret, err
	}
	return ret, nil
}
