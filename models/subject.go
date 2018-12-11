package models

import (
	"sort"

	"github.com/astaxie/beego/logs"
)

// Sm is global subject manager
var Sm SubjectManager

// SubjectManager manager all course subject
type SubjectManager struct {
	subject map[int]SubjectInfo
}

// Add add new subject into manager
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

// Update update the key of subject
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

// Delete remove subject of id
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
func (sm *SubjectManager) GetAll() SubjectList {
	ret := SubjectList{}
	for k, v := range sm.subject {
		ret = append(ret, SubjectInfo{ID: k, Name: v.Name, Key: v.Key})
	}

	sort.Sort(ret)
	return ret
}

// CheckSubjectList check to see if all id in input list exist
func (sm *SubjectManager) CheckSubjectList(list []int) bool {
	for _, v := range list {
		if _, ok := sm.subject[v]; !ok {
			return false
		}
	}
	return true
}

// IsExist check to see if exist
func (sm *SubjectManager) IsExist(id int) bool {
	_, ok := sm.subject[id]
	return ok
}

func (sm *SubjectManager) getSubjectName(id int) string {
	return sm.subject[id].Name
}
