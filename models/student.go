package models

import (
	"time"
	"unicode/utf8"

	"github.com/arong/dean/base"
	"github.com/bearbin/go-age"
)

type StudentInfo struct {
	Age        int    `json:"age"`
	Gender     int    `json:"gender"`
	Name       string `json:"name"` // required
	Mobile     string `json:"mobile"`
	Address    string `json:"address"`
	Birthday   string `json:"birthday"`
	StudentID  int64  `json:"student_id"`
	Status     int    `json:"status"`
	ClassID    int    `json:"class_id"`    // required
	RegisterID string `json:"register_id"` // required
}

func (s StudentInfo) Equal(r StudentInfo) bool {
	return s.Gender == r.Gender &&
		s.Name == r.Name &&
		s.Mobile == r.Mobile &&
		s.Address == r.Address &&
		s.Birthday == r.Birthday &&
		s.ClassID == r.ClassID &&
		s.RegisterID == r.RegisterID
}

func (s StudentInfo) Check() error {
	if s.Gender < eGenderMale || s.Gender > eGenderUnknown {
		return ErrGender
	}

	if len(s.Name) == 0 {
		return ErrName
	}

	if len(s.Mobile) > 11 {
		return errMobile
	}

	if utf8.RuneCountInString(s.Address) > 64 {
		return errAddress
	}

	if len(s.Birthday) > 10 {
		return ErrBirthday
	}

	if len(s.Birthday) != 0 {
		tmp, err := time.Parse(base.DateFormat, s.Birthday)
		if err != nil {
			return ErrBirthday
		}
		s.Age = age.Age(tmp)
	}

	if len(s.RegisterID) == 0 {
		return errRegisterID
	}

	for _, v := range s.RegisterID {
		if v > '9' || v < '0' {
			return errRegisterID
		}
	}
	return nil
}

type StudentList []StudentInfo

func (s StudentList) Filter(f StudentFilter) (int, StudentList) {
	ret := StudentList{}
	start, end := f.GetRange()

	i := 0
	total := 0
	for _, v := range s {
		if f.Name != "" && f.Name != v.Name {
			continue
		}

		if f.RegisterID != "" && f.RegisterID != v.RegisterID {
			continue
		}

		// satisfy requirement
		total++
		if i >= start && i < end {
			ret = append(ret, v)
		}
		i++
	}

	return total, ret
}

func (cl StudentList) Len() int {
	return len(cl)
}

func (cl StudentList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

func (cl StudentList) Less(i, j int) bool {
	return cl[i].StudentID < cl[j].StudentID
}

//go:generate mockgen -destination=../mocks/mock_student.go -package mocks github.com/arong/dean/models StudentStore
type StudentStore interface {
	SaveStudent(StudentInfo) (int64, error)
	UpdateStudent(info StudentInfo) error
	DeleteStudent([]int64) error
}
