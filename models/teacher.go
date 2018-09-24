package models

import "errors"

type Teacher struct {
	ID     int
	Gender int    // 1: Male, 2: Female, 3: unknown
	Name   string // 姓名
	Mobile string // 手机号
}

type TeacherManager struct {
	nameMap map[string]*Teacher // 名字索引
	idMap   map[int]*Teacher    // id索引表
}

var (
	ErrInvalidParam = errors.New("invalid parameter")
	ErrNameExist    = errors.New("name exist")
	ErrNotExist     = errors.New("teacher not exist")
)

// a global handler
var tm TeacherManager

func (tm *TeacherManager) AddTeacher(t *Teacher) error {
	if t.ID != 0 {
		return ErrInvalidParam
	}

	if _, ok := tm.nameMap[t.Name]; ok {
		return ErrNameExist
	}
	// todo: validate parameter

	// add to map
	tm.nameMap[t.Name] = t

	return nil
}

func (tm *TeacherManager) DelTeacher(id int) error {
	val, ok := tm.idMap[id]
	if !ok {
		return ErrNotExist
	}

	// remove from map
	delete(tm.nameMap, val.Name)
	delete(tm.idMap, val.ID)

	return nil
}

func (tm *TeacherManager) GetTeacherInfo(id int) (*Teacher, error) {
	ret := &Teacher{}
	err := ErrNotExist
	if val, ok := tm.idMap[id]; ok {
		*ret = *val
		err = nil
	}
	return ret, err
}

func (tm *TeacherManager) GetTeacherList(ids []int) ([]*Teacher, error) {
	ret := []*Teacher{}
	for _, v := range ids {
		if val, ok := tm.idMap[v]; ok {
			tmp := *val
			ret = append(ret, &tmp)
		}
	}
	return ret, nil
}
