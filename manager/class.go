package manager

import (
	"errors"
	"sort"
	"sync"

	"github.com/arong/dean/models"

	"github.com/arong/dean/base"
	"github.com/astaxie/beego/logs"
)

// global manager
var Cm classManager

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

type classManager struct {
	idMap map[int]*models.Class
	mutex sync.Mutex
}

// Init maintain the relation between class and teacher
func (cm *classManager) Init(data map[int]*models.Class) {
	if cm == nil {
		return
	}

	if data == nil {
		cm.idMap = make(map[int]*models.Class)
	} else {
		cm.idMap = data
	}
}

// AddClass add new class into system
func (cm *classManager) AddClass(c *models.Class) (int, error) {
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

//ModifyClass modify class
func (cm *classManager) ModifyClass(r *models.Class) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	curr, ok := cm.idMap[r.ID]
	if !ok {
		return ErrClassNotExist
	}

	// todo: fix up this
	//err := Tm.CheckInstructorList(r.TeacherList)
	//if err != nil {
	//	logs.Warn("[ModifyClass] CheckInstructorList error", "err", err)
	//	return err
	//}

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
	err := Ma.UpdateClass(curr)
	if err != nil {
		logs.Warn("[ModifyClass] database error")
		return err
	}
	curr.AddList = models.InstructorList{}
	curr.RemoveList = models.InstructorList{}
	return nil
}

//DelClass delete class
func (cm *classManager) DelClass(list models.ClassIDList) (models.ClassIDList, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	failedList := models.ClassIDList{}

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

// Filter  get class list with condition
func (cm *classManager) Filter() base.CommList {
	resp := base.CommList{}
	list := models.ClassList{}
	for _, v := range cm.idMap {
		list = append(list, v)
	}
	sort.Sort(list)
	resp.List = list
	resp.Total = len(list)
	return resp
}

// GetAll get all class list
func (cm *classManager) GetAll() models.ItemList {
	resp := models.ItemList{}
	for _, v := range cm.idMap {
		resp = append(resp, models.Item{ID: v.ID, Name: v.Name})
	}
	sort.Sort(resp)
	return resp
}

// GetInfo get class info
func (cm *classManager) GetInfo(id int) (*models.Class, error) {
	ret := &models.Class{}

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
