package models

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bearbin/go-age"

	"github.com/arong/dean/base"
)

const (
	eGenderMale    = 1
	eGenderFemale  = 2
	eGenderUnknown = 3
)

type TeacherMeta struct {
	TeacherID int64  `json:"teacher_id"`
	SubjectID int    `json:"subject_id"`
	Gender    int    `json:"gender"`
	Age       int    `json:"age"`
	Name      string `json:"name,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Birthday  string `json:"birthday,omitempty"`
	Address   string `json:"address,omitempty"`
	Subject   string `json:"subject,omitempty"`
}

func (t TeacherMeta) Equal(r TeacherMeta) bool {
	return t.SubjectID == r.SubjectID &&
		t.Gender == r.Gender &&
		t.Name == r.Name &&
		t.Mobile == r.Mobile &&
		t.Birthday == r.Birthday &&
		t.Address == r.Address &&
		t.Subject == r.Subject
}

type Teacher struct {
	Status int `json:"-"`
	TeacherMeta
}

func (t Teacher) Check() error {
	if t.Gender < eGenderMale || t.Gender > eGenderUnknown {
		return ErrGender
	}

	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return ErrName
	}

	if utf8.RuneCountInString(t.Name) > 16 {
		return ErrName
	}

	if len(t.Mobile) > 11 {
		return errMobile
	}

	if len(t.Birthday) > 10 {
		return ErrBirthday
	}

	if t.Birthday != "" {
		birth, err := time.Parse("2006-01-02", t.Birthday)
		if err != nil {
			return ErrBirthday
		}
		t.Age = age.Age(birth)
	}

	if utf8.RuneCountInString(t.Address) > 64 {
		return errAddress
	}

	// todo: fixup
	//if t.SubjectID != 0 {
	//	if !Sm.IsExist(t.SubjectID) {
	//		return errSubject
	//	}
	//}
	return nil
}

type TeacherList []Teacher

func (tl TeacherList) Len() int {
	return len(tl)
}

func (tl TeacherList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (tl TeacherList) Less(i, j int) bool {
	return tl[i].TeacherID < tl[j].TeacherID
}

func (tl TeacherList) Page(page, size int) TeacherList {
	start := (page - 1) * size
	end := page * size
	total := len(tl)

	if start >= total {
		return TeacherList{}
	} else if end > total {
		return tl[start:]
	} else {
		return tl[start:end]
	}
}

type TeacherFilter struct {
	base.CommPage
	Gender int    `json:"gender"`
	Age    int    `json:"age"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

func (tl TeacherList) Filter(f TeacherFilter) (int, TeacherList) {
	list := TeacherList{}
	start, end := f.GetRange()

	i := 0
	total := 0
	for _, v := range tl {
		if v.Status != base.StatusValid {
			continue
		}

		if f.Gender != 0 && f.Gender != v.Gender {
			continue
		}

		if f.Name != "" && f.Name != v.Name {
			continue
		}

		if f.Mobile != "" && f.Mobile != v.Mobile {
			continue
		}

		if i >= start && i < end {
			list = append(list, v)
		}
		i++
		total++
	}
	return total, list
}

//go:generate mockgen -destination=../mocks./mock_teacher.go -package mocks github.com/arong/dean/models TeacherStore
type TeacherStore interface {
	SaveTeacher(Teacher) (int64, error)
	UpdateTeacher(Teacher) error
	DeleteTeacher([]int64) error
}
