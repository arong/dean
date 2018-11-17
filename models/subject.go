package models

import "sort"

// Sm is global subject manager
var Sm SubjectManager

// SubjectManager manager all course subject
type SubjectManager struct {
	subject map[int]string
}

type SubjectInfo struct {
	SubjectID   int    `json:"subject_id"`
	SubjectName string `json:"subject_name"`
}

type SubjectList []*SubjectInfo

func (tl SubjectList) Len() int {
	return len(tl)
}

func (tl SubjectList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (tl SubjectList) Less(i, j int) bool {
	return tl[i].SubjectID < tl[j].SubjectID
}

// GetAll return all subject id list in current manager
func (s *SubjectManager) GetAll() SubjectList {
	ret := SubjectList{}
	for k, v := range s.subject {
		ret = append(ret, &SubjectInfo{SubjectID: k, SubjectName: v})
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
