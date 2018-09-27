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
	Grade int // 年级
	Index int // 班级
}

func (f *Filter) GetID() ClassID {
	if f == nil {
		return 0
	}
	return ClassID((f.Grade & 0xf0 << 8) | f.Index&0x0f)
}

type ClassManager struct {
	idMap map[ClassID]*Class
}

var (
	ErrClassNotExist = errors.New("class not exist")
)

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

func (cm *ClassManager) GetTeacherList(grade, index int) (TeacherList, error) {
	ret := TeacherList{}
	key := Filter{
		Grade: grade,
		Index: index,
	}

	val, ok := cm.idMap[key.GetID()]
	if !ok {
		return ret, ErrClassNotExist
	}
	return Tm.GetTeacherList(val.TeacherIDs)
}
