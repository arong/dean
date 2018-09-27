package models

import (
	"errors"
)

// global manager
var Cm ClassManager
type ClassID int
type Class struct {
	Filter
	ID         ClassID
	Name       string
	TeacherIDs []int64 // 教师id
}

type ClassList []*Class

type Filter struct {
	Grade      int     // 年级
	Index      int     // 班级
}
type ClassManager struct {
	idMap map[Filter]*Class
}

var (
	ErrClassNotExist = errors.New("class not exist")
)

// maintain the relation between class and teacher
func (cm *ClassManager) Init() {
	cm.idMap = make(map[Filter]*Class)
}

func (cm *ClassManager) GetTeacherList(grade, index int) (TeacherList, error) {
	ret := TeacherList{}
	key := Filter{
		Grade: grade,
		Index: index,
	}

	val, ok := cm.idMap[key]
	if !ok {
		return ret, ErrClassNotExist
	}
	return Tm.GetTeacherList(val.TeacherIDs)
}

