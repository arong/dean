package models

import (
	"errors"
)

// global manager
var Cm ClassManager

// of size two byte
// +--------+-------+
// | Grade  | index |
// +--------+-------+
type ClassID int

// Class is the
type Class struct {
	Filter
	ID         ClassID
	Name       string
	TeacherIDs []int64 `json:"-"` // id
}

type ClassResp struct {
	Class
	Teachers TeacherList // 详情
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
	return ClassID(((f.Grade & 0xf) << 8) | f.Index&0x0f)
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

func (cm *ClassManager) GetAll() ClassList {
	ret := ClassList{}
	for _, v := range cm.idMap {
		ret = append(ret, v)
	}
	return ret
}

func (cm *ClassManager) GetInfo(f *Filter) (*ClassResp, error) {
	ret := &ClassResp{}
	var err error

	if cm == nil {
		return ret, nil
	}

	if f == nil {
		return ret, errors.New("invalid input")
	}

	val, ok := cm.idMap[f.GetID()]
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
