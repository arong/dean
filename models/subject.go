package models

import "sort"

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

func (s *SubjectManager) GetAll() SubjectList {
	ret := SubjectList{}
	for k, v := range s.subject {
		ret = append(ret, &SubjectInfo{SubjectID: k, SubjectName: v})
	}

	sort.Sort(ret)
	return ret
}
