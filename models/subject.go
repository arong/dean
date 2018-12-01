package models

import (
	"github.com/astaxie/beego/logs"
	"sort"
)

// Sm is global subject manager
var Sm SubjectManager

// SubjectManager manager all course subject
type SubjectManager struct {
	subject map[int]SubjectInfo
}

type SubjectInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
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

func (sm *SubjectManager) Add(s SubjectInfo) (int, error) {
	var err error
	for _, v := range sm.subject {
		if s.Key == v.Key || s.Name == v.Name {
			logs.Info("[SubjectManager::Add] conflict", "key", s.Key, "name", s.Name)
			return 0, errExist
		}
	}

	s.ID, err = Ma.AddSubject(&s)
	if err != nil {
		logs.Info("[SubjectManager::Add] db failed", "err", err)
		return 0, err
	}
	// add to current map
	sm.subject[s.ID] = s
	return s.ID, nil
}

func (sm *SubjectManager) Update(s SubjectInfo) error {
	var err error
	curr, ok := sm.subject[s.ID]
	if !ok {
		logs.Info("[SubjectManager::Update] not found", "id", s.ID)
		return errNotExist
	}

	if s.Key == curr.Key {
		logs.Info("[SubjectManager::Update] no change")
		return nil
	}

	for _, v := range sm.subject {
		if s.Key == v.Key && s.ID != v.ID {
			logs.Info("[SubjectManager::Update] conflict", "key", s.Key, "name", s.Name)
			return errExist
		}
	}

	err = Ma.UpdateSubject(&s)
	if err != nil {
		logs.Info("[SubjectManager::Update] db failed", "err", err)
		return err
	}

	// update cache
	curr.Key = s.Key
	sm.subject[s.ID] = s

	return nil
}

func (sm *SubjectManager) Delete(id int) error {
	var err error
	if _, ok := sm.subject[id]; !ok {
		logs.Info("[SubjectManager::Delete] class not found", "id", id)
		return errNotExist
	}

	err = Ma.DeleteSubject(id)
	if err != nil {
		logs.Info("[SubjectManager::Delete] db failed", "err", err)
		return err
	}
	// add to current map
	delete(sm.subject, id)
	return nil
}

// GetAll return all subject id list in current manager
func (s *SubjectManager) GetAll() SubjectList {
	ret := SubjectList{}
	for k, v := range s.subject {
		ret = append(ret, SubjectInfo{ID: k, Name: v.Name, Key: v.Key})
	}

	sort.Sort(ret)
	return ret
}

// CheckSubjectList check to see if all id in input list exist
func (s *SubjectManager) CheckSubjectList(list []int) bool {
	for _, v := range list {
		if _, ok := s.subject[v]; !ok {
			return false
		}
	}
	return true
}
