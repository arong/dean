package models

import (
	"github.com/pkg/errors"
	"sort"
)

/*
 * struct info
 */

// ClassID of size two byte
// +--------+-------+
// | Grade  | index |
// +--------+-------+
type ClassID int

type ClassIDList []int

// Class is the
type Class struct {
	Filter

	ID          int            `json:"id"`
	MasterID    UserID         `json:"master_id"` // 班主任
	Season      int            `json:"season"`    // 1: 春季, 3: 秋季
	Name        string         `json:"name"`      // 班级名称
	Year        int            `json:"year"`      // 所在年份
	TeacherList InstructorList `json:"teacher_list"`
	RemoveList  InstructorList `json:"-"`
	AddList     InstructorList `json:"-"`
}

func (c Class) Check() error {
	if c.Grade <= 0 {
		return errors.New("invalid grade")
	}

	if c.Index <= 0 {
		return errors.New("invalid index")
	}

	if c.Year == 0 {
		return errors.New("invalid year")
	}

	return nil
}

func (c Class) Equal(r Class) bool {
	if c.MasterID != r.MasterID ||
		c.Name != r.Name ||
		c.Year != r.Year ||
		c.Season != r.Season {
		return false
	}
	if len(c.TeacherList) != len(r.TeacherList) {
		return false
	}
	sort.Sort(c.TeacherList)
	sort.Sort(r.TeacherList)
	for k, v := range c.TeacherList {
		if v.SubjectID != r.TeacherList[k].SubjectID ||
			v.TeacherID != r.TeacherList[k].TeacherID {
			return false
		}
	}
	return true
}

type ClassList []*Class

type ClassResp struct {
	Class
	Teachers TeacherList // 详情
}

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

//func (f *Filter) GetID() ClassID {
//	if f == nil {
//		return 0
//	}
//	return ClassID(((f.Grade & 0xf) << 8) | f.Index&0x0f)
//}

// InstructorInfo specify teacher and its subject id
type InstructorInfo struct {
	TeacherID UserID `json:"teacher_id"`
	SubjectID int    `json:"subject_id"`
}

type InstructorList []InstructorInfo

func (il InstructorList) Len() int {
	return len(il)
}

func (il InstructorList) Swap(i, j int) {
	il[j], il[i] = il[i], il[j]
}

func (il InstructorList) Less(i, j int) bool {
	if il[i].SubjectID < il[j].SubjectID {
		return true
	} else if il[i].SubjectID > il[j].SubjectID {
		return false
	} else {
		return il[i].TeacherID < il[j].TeacherID
	}
}
