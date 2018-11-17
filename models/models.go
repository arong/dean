package models

/*
 * struct info
 */

// ClassID of size two byte
// +--------+-------+
// | Grade  | index |
// +--------+-------+
type ClassID int

type ClassIDList []ClassID

// Class is the
type Class struct {
	Filter
	ID             ClassID `json:"id"`
	Name           string  `json:"name"`
	MasterID       UserID  `json:"master_id"` // 班主任
	InstructorList []InstructorInfo
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

func (f *Filter) GetID() ClassID {
	if f == nil {
		return 0
	}
	return ClassID(((f.Grade & 0xf) << 8) | f.Index&0x0f)
}

// InstructorInfo specify teacher and its subject id
type InstructorInfo struct {
	TeacherID UserID `json:"teacher_id"`
	SubjectID int    `json:"subject_id"`
}

type InstructorList []InstructorInfo
