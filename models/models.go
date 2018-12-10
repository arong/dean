package models

import (
	"github.com/arong/dean/base"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"sort"
	"time"
)

/*
 * struct info
 */

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ItemList []Item

func (il ItemList) Len() int {
	return len(il)
}
func (il ItemList) Swap(i, j int) {
	il[j], il[i] = il[i], il[j]
}
func (il ItemList) Less(i, j int) bool {
	return il[i].ID < il[j].ID
}

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
	Term        int            `json:"term"`      // 1: 第一学期, 3: 第二学期
	Name        string         `json:"name"`      // 班级名称
	Year        int            `json:"year"`      // 所在年份
	TeacherList InstructorList `json:"teacher_list,omitempty"`
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

	if c.Term == 0 {
		return errors.New("invalid season")
	}

	if c.MasterID == 0 {
		return errors.New("invalid master id")
	}
	return nil
}

func (c Class) Equal(r Class) bool {
	if c.MasterID != r.MasterID ||
		c.Name != r.Name ||
		c.Year != r.Year ||
		c.Term != r.Term {
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

type StudentFilter struct {
	Filter
	base.CommPage
	Name   string
	Number string
}

func (s StudentFilter) Check() error {
	if s.Name != "" {
		return nil
	}

	if s.Number != "" {
		return nil
	}

	if s.Page < 0 || s.Size <= 0 {
		return errors.New("page error")
	}
	return nil
}

// InstructorInfo specify teacher and its subject id
type InstructorInfo struct {
	TeacherID UserID `json:"tid"`
	SubjectID int    `json:"sid"`
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

func (il InstructorList) Deduplicate() InstructorList {
	tmp := make(map[int]InstructorInfo)
	for _, v := range il {
		val := v
		if v.SubjectID == 0 || v.TeacherID == 0 {
			continue
		}
		tmp[v.SubjectID] = val
	}

	// no duplicated item
	if len(il) == len(tmp) {
		return il
	}

	newList := make([]InstructorInfo, 0, len(tmp))
	for _, v := range tmp {
		val := v
		newList = append(newList, val)
	}
	logs.Info("[]", "newList", newList)
	return newList
}

func (il InstructorList) Diff(r InstructorList) (all, add, del InstructorList) {
	all, add, del = InstructorList{}, InstructorList{}, InstructorList{}
	curr := make(map[int]InstructorInfo)
	for _, v := range il {
		curr[v.SubjectID] = v
	}

	for _, i := range r {
		v := i
		if val, ok := curr[v.SubjectID]; ok {
			tmp := val
			if v.TeacherID != val.TeacherID {
				add = append(add, v)
				del = append(del, tmp)
			}
			delete(curr, v.SubjectID)
		} else {
			add = append(add, v)
		}
	}

	for _, i := range il {
		v := i
		if _, ok := curr[v.SubjectID]; ok {
			continue
		}
		all = append(all, v)
	}

	all = append(all, add...)

	// delete list
	for _, v := range curr {
		tmp := v
		del = append(del, tmp)
	}

	return all, add, del
}

// questionnaire zone
type OptionInfo struct {
	OptionID int    `json:"option_id"`
	Index    int    `json:"index"`
	Option   string `json:"option"`
}

func (o OptionInfo) Check() error {
	if o.Option == "" {
		return errors.New("invalid option")
	}
	if o.Index < 0 {
		return errors.New("invalid index")
	}
	return nil
}

type OptionList []OptionInfo

func (ol OptionList) Len() int {
	return len(ol)
}
func (ol OptionList) Swap(i, j int) {
	ol[j], ol[i] = ol[i], ol[j]
}
func (ol OptionList) Less(i, j int) bool {
	return ol[i].Index < ol[j].Index
}

func (ol OptionList) Check() error {
	option := make(map[string]bool)
	index := make(map[int]bool)
	for _, v := range ol {
		err := v.Check()
		if err != nil {
			return err
		}
		if _, ok := option[v.Option]; ok {
			return errors.New("option duplicate")
		} else {
			option[v.Option] = true
		}
		if _, ok := index[v.Index]; ok {
			return errors.New("index duplicate")
		} else {
			index[v.Index] = true
		}
	}
	return nil
}

type QuestionInfo struct {
	QuestionID int        `json:"question_id"`
	Index      int        `json:"index"`
	Type       int        `json:"type"`
	Required   bool       `json:"required"`
	Question   string     `json:"question"`
	Options    OptionList `json:"options"`
	Scope      []int      `json:"scope"` // which subject will this question apply to
}

func (q QuestionInfo) Check() error {
	if q.Question == "" {
		return errors.New("invalid question")
	}

	if q.Type < QuestionTypeSingleSelection || q.Type > QuestionTypeText {
		return errors.New("invalid question type")
	}

	if len(q.Options) == 0 {
		return errors.New("empty questions")
	}
	err := q.Options.Check()
	if err != nil {
		return err
	}

	if len(q.Scope) > 0 {
		if !Sm.CheckSubjectList(q.Scope) {
			return errors.New("invalid scope")
		}
	}
	return err
}

type QuestionList []QuestionInfo

func (q QuestionList) Len() int {
	return len(q)
}
func (q QuestionList) Swap(i, j int) {
	q[j], q[i] = q[i], q[j]
}
func (q QuestionList) Less(i, j int) bool {
	return q[i].Index < q[j].Index
}
func (q QuestionList) Check() error {
	title := make(map[string]bool)
	index := make(map[int]bool)

	for _, v := range q {
		err := v.Check()
		if err != nil {
			return err
		}

		if _, ok := title[v.Question]; ok {
			return errors.New("duplicate question found")
		} else {
			title[v.Question] = true
		}

		if _, ok := index[v.Index]; ok {
			return errors.New("index duplicated")
		} else {
			index[v.Index] = true
		}
	}
	return nil
}

type QuestionnaireInfo struct {
	QuestionnaireID int          `json:"id"`
	Status          int          `json:"status"` // 1: draft, 2: published
	Title           string       `json:"title"`
	StartTime       string       `json:"start"`
	StopTime        string       `json:"stop"`
	Questions       QuestionList `json:"questions"`
	startTime       time.Time
	stopTime        time.Time
}

type QuestionnaireList []QuestionnaireInfo

func (ql QuestionnaireList) Len() int {
	return len(ql)
}
func (ql QuestionnaireList) Swap(i, j int) {
	ql[j], ql[i] = ql[i], ql[j]
}
func (ql QuestionnaireList) Less(i, j int) bool {
	return ql[i].startTime.Before(ql[j].startTime)
}
func (ql QuestionnaireList) Page(page base.CommPage) QuestionnaireList {
	start, end := page.GetRange()
	size := len(ql)
	ret := QuestionnaireList{}
	if start >= size {
		return ret
	} else if end > size {
		return ql[start:]
	} else {
		return ql[start:end]
	}
}

const (
	QuestionTypeSingleSelection = 1
	QuestionTypeMultiSelection  = 2
	QuestionTypeText            = 3
	// questionnaire status
	QStatusDraft     = 1 // draft, could be edited
	QStatusPublished = 2 // published, could not be edited
	QStatusDrawBack  = 3
	QStatusExpired   = 4
)

func (q QuestionnaireInfo) Check() error {
	if q.Title == "" {
		return errors.New("invalid title")
	}

	if q.Status != QStatusDraft && q.Status != QStatusPublished {
		return errors.New("invalid draft status")
	}

	if q.StopTime != "" && q.stopTime.Before(time.Now()) {
		return errors.New("invalid stop time")
	}

	if q.StartTime != "" && q.StopTime != "" && q.startTime.After(q.stopTime) {
		return errors.New("invalid time range")
	}

	if len(q.Questions) == 0 {
		return errors.New("empty questions")
	}
	err := q.Questions.Check()

	return err
}
