package models

import (
	"errors"
)

type Class struct {
	ID         int
	Name       string
	Grade      int   // 年级
	Index      int   // 班级
	TeacherIDs []int // 教师id
}

type Filter struct {
	Grade int
	Index int
}
type ClassManager struct {
	idMap map[Filter]*Class
}

var (
	ErrClassNotExist = errors.New("class not exist")
)

// maintain the relation between class and teacher
func (cm *ClassManager) GetTeacherList(grade, index int) ([]*Teacher, error) {
	ret := []*Teacher{}
	key := Filter{
		Grade: grade,
		Index: index,
	}

	val, ok := cm.idMap[key]
	if !ok {
		return ret, ErrClassNotExist
	}
	return tm.GetTeacherList(val.TeacherIDs)
}
