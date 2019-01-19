package models

import "github.com/pkg/errors"

// SubjectInfo: subject meta info
type SubjectInfo struct {
	Status int    `json:"status"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
}

func (s SubjectInfo) Check() error {
	for _, v := range s.Key {
		if (v >= 'a' && v <= 'z') ||
			(v >= 'A' && v <= 'Z') ||
			v == '_' {
			continue
		}
		return errors.New("invalid key")
	}
	return nil
}

func (s SubjectInfo) Equal(r SubjectInfo) bool {
	return s.Name == r.Name &&
		s.Key == r.Key
}

type SubjectList []SubjectInfo

func (tl SubjectList) Len() int {
	return len(tl)
}

func (tl SubjectList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (tl SubjectList) Less(i, j int) bool {
	return tl[i].ID < tl[j].ID
}

//go:generate mockgen -destination=../mocks/mock_subject.go -package mocks github.com/arong/dean/models SubjectStore
type SubjectStore interface {
	SaveSubject(SubjectInfo) (int, error)
	UpdateSubject(SubjectInfo) error
	DeleteSubject(int) error
}
