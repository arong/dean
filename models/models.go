package models

import (
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bearbin/go-age"

	"github.com/astaxie/beego"

	"github.com/arong/dean/base"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
)

var (
	ErrName     = errors.New("invalid name")
	ErrGender   = errors.New("invalid gender")
	ErrBirthday = errors.New("invalid birthday")
	errAddress  = errors.New("invalid address")
	errMobile   = errors.New("invalid mobile number")
	errSubject  = errors.New("invalid subject id")
)

/*
 * struct info
 */

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"DBName"`
}

// GetConf get config from app.conf
func (c *DBConfig) GetConf() error {
	c.User = beego.AppConfig.String("user")
	c.Password = beego.AppConfig.String("password")
	c.Host = beego.AppConfig.String("host")
	c.Port, _ = beego.AppConfig.Int("port")
	c.DBName = beego.AppConfig.String("dbName")

	logs.Debug(*c)
	return nil
}

type IntList []int

func (il IntList) Page(page base.CommPage) IntList {
	start, end := page.GetRange()
	size := len(il)
	ret := IntList{}
	if start >= size {
		return ret
	} else if end > size {
		return il[start:]
	} else {
		return il[start:end]
	}
}

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
	Term        int            `json:"term"`      // 1: 第一学期, 3: 第二学期
	MasterID    int64          `json:"master_id"` // 班主任
	Name        string         `json:"name"`      // 班级名称
	Year        int            `json:"year"`      // 所在年份
	TeacherList InstructorList `json:"teacher_list,omitempty"`
	RemoveList  InstructorList `json:"-"`
	AddList     InstructorList `json:"-"`
	StudentList []int64        `json:"-"`
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
	TeacherID int64  `json:"tid"`
	SubjectID int    `json:"sid"`
	Subject   string `json:"subject"`
	Teacher   string `json:"teacher"`
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
	return newList
}

// Diff: diff request with current list, and get merged, newly added, deleted list
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
	Index  int    `json:"index"`
	Option string `json:"option"`
}

func (o OptionInfo) Equal(r OptionInfo) bool {
	return o.Index == r.Index && o.Option == r.Option
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

func (ol OptionList) FilterEmpty() OptionList {
	ret := OptionList{}
	for _, v := range ol {
		if v.Index == 0 || v.Option == "" {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

type QCheckBox struct {
	Choice OptionList
}

type QScore struct {
	Min int
	Max int
}

type QStar struct {
	Count int
}

type QuestionInfo struct {
	QuestionnaireID int        `json:"questionnaire_id"`
	QuestionID      int        `json:"id"`
	Index           int        `json:"index"`
	Type            int        `json:"type"`
	Required        bool       `json:"required"`
	Question        string     `json:"question"`
	Options         OptionList `json:"options"`
	Scope           []int      `json:"scope"` // which subject will this question apply to
}

func (q QuestionInfo) Equal(r *QuestionInfo) bool {
	if q.QuestionID != r.QuestionID ||
		q.Index != r.Index ||
		q.Type != r.Type ||
		q.Required != r.Required ||
		q.Question != r.Question ||
		len(q.Options) != len(r.Options) ||
		len(q.Scope) != len(r.Scope) {
		return false
	}

	for k, v := range q.Scope {
		if r.Scope[k] != v {
			return false
		}
	}

	sort.Sort(q.Options)
	sort.Sort(r.Options)

	for k, v := range q.Options {
		if !r.Options[k].Equal(v) {
			return false
		}
	}
	return true
}

func (q QuestionInfo) Check() error {
	if q.Question == "" {
		return errors.New("invalid question")
	}

	if q.Type < QuestionTypeSingleSelection || q.Type > QuestionTypeText {
		return errors.New("invalid question type")
	}

	if len(q.Options) == 0 {
		return errors.New("empty options")
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

type QuestionList []*QuestionInfo

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
	Label           string       `json:"label"`
	Questions       QuestionList `json:"questions"`
	Editor          string       `json:"editor"`
	startTime       time.Time
	stopTime        time.Time
}

func (q *QuestionnaireInfo) Check() error {
	if q.Title == "" {
		return errors.New("invalid title")
	}

	if q.Status != QStatusDraft && q.Status != QStatusPublished {
		return errors.New("invalid draft status")
	}

	var err error
	if q.StartTime != "" {
		q.startTime, err = time.Parse(base.DateTimeFormat, q.StartTime)
		if err != nil {
			return errors.New("invalid start time format")
		}
	}

	if q.StopTime != "" {
		q.stopTime, err = time.Parse(base.DateTimeFormat, q.StopTime)
		if err != nil {
			return errors.New("invalid stop time format")
		}
	} else {
		return errors.New("stop time needed")
	}

	if q.stopTime.Before(time.Now()) {
		return errors.New("invalid stop time")
	}

	if q.startTime.After(q.stopTime) {
		return errors.New("invalid time range")
	}

	return err
}

func (q QuestionnaireInfo) Equal(r QuestionnaireInfo) bool {
	if q.Status != r.Status ||
		q.Title != r.Title ||
		q.StartTime != r.StartTime ||
		q.StopTime != r.StopTime {
		return false
	}
	return true
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

type SurveyPage struct {
	TeacherID   int64        `json:"t_id"`
	TeacherName string       `json:"t_name"`
	Questions   QuestionList `json:"questions"`
}

type SurveyPages []SurveyPage

func (s SurveyPages) Len() int {
	return len(s)
}
func (s SurveyPages) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SurveyPages) Less(i, j int) bool {
	return s[i].TeacherID < s[j].TeacherID
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

type AnswerInfo struct {
	QuestionID int
	Answer     interface{}
}
type AnswerList []*AnswerInfo

func (al AnswerList) Len() int {
	return len(al)
}
func (al AnswerList) Swap(i, j int) {
	al[i], al[j] = al[j], al[i]
}
func (al AnswerList) Less(i, j int) bool {
	return al[i].QuestionID < al[j].QuestionID
}
func (al AnswerList) Check() error {
	tmp := make(map[int]bool)
	for _, v := range al {
		if v.QuestionID == 0 {
			return errNotExist
		}

		if _, ok := tmp[v.QuestionID]; ok {
			return errExist
		} else {
			tmp[v.QuestionID] = true
		}
	}
	return nil
}

type TeacherAnswer struct {
	TeacherID int64
	Answers   AnswerList
}
type TeacherAnswerList []TeacherAnswer

type QuestionnaireSubmit struct {
	QuestionnaireID int
	StudentID       int64
	TeacherAnswers  TeacherAnswerList
}

func (q QuestionnaireSubmit) Check() error {
	if q.QuestionnaireID == 0 {
		return errNotExist
	}
	return nil
}

// ScorePair is score of single subject
type ScorePair struct {
	SubjectID int
	Score     int
}

type ScorePairList []ScorePair

type ExamScore struct {
	Exam   int
	Scores ScorePairList
}

type ExamScoreList []ExamScore

type TermScore struct {
	TermID      int
	ExamsScores ExamScoreList
}

type TermScoreList []TermScore

type YearScore struct {
	Year       int
	TermScores [2]TermScore
}

type YearScoreList []YearScore

// StudentScore request of adding score record
type StudentScore struct {
	StudentID int64
	TermID    int
	Exam      int
	Scores    ScorePairList
}

func (ss StudentScore) Check() error {
	if ss.StudentID == 0 {
		return errors.New("invalid student info")
	} else {
		if !Um.IsExist(ss.StudentID) {
			return errors.New("student not exist")
		}
	}

	if ss.TermID == 0 || ss.Exam == 0 {
		return errors.New("invalid score")
	}

	for _, v := range ss.Scores {
		if v.SubjectID == 0 {
			return errors.New("invalid subject id")
		}
		// check subject
		if !Sm.IsExist(v.SubjectID) {
			return errNotExist
		}

		if v.Score < base.MinScore || v.Score > base.MaxScore {
			return errors.New("invalid score")
		}
	}
	return nil
}

type StudentScoreList []StudentScore

func (ssl StudentScoreList) Len() int {
	return len(ssl)
}

func (ssl StudentScoreList) Swap(i, j int) {
	ssl[i], ssl[j] = ssl[j], ssl[i]
}

func (ssl StudentScoreList) Less(i, j int) bool {
	return ssl[i].StudentID < ssl[j].StudentID
}

type ScoreFilter struct {
	OnlyCurrent bool  // 是否只是当前学期
	SubjectID   int   // 课程ID
	ClassID     int   // 班级
	StudentID   int64 // 学生ID
	TermID      int   // 学期
}

type StudentInfo struct {
	profile
	ClassID    int    `json:"class_id"`
	StudentID  int64  `json:"student_id"`
	RegisterID string `json:"register_id"` // 学号
}
type studentList []*StudentInfo

func (cl studentList) Len() int {
	return len(cl)
}

func (cl studentList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

func (cl studentList) Less(i, j int) bool {
	return cl[i].StudentID < cl[j].StudentID
}

type profile struct {
	Age      int    `json:"age"`
	Gender   int    `json:"gender"`
	RealName string `json:"real_name"`
	Mobile   string `json:"mobile"`
	Address  string `json:"address"`
	Birthday string `json:"birthday"`
}

// SubjectInfo: subject meta info
type SubjectInfo struct {
	Status int    `json:"status"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
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

	if t.SubjectID != 0 {
		if !Sm.IsExist(t.SubjectID) {
			return errSubject
		}
	}
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

func (tl TeacherList) Filter(f TeacherFilter) TeacherList {
	list := TeacherList{}
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
		list = append(list, v)
	}
	return list
}

type TeacherFilter struct {
	base.CommPage
	Gender int    `json:"gender"`
	Age    int    `json:"age"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}

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

type TeacherInfoResp struct {
	TeacherMeta
	SubjectID int `json:"subject_id"`
}

type simpleTeacher struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type simpleTeacherList []simpleTeacher

func (s simpleTeacherList) Len() int {
	return len(s)
}
func (s simpleTeacherList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s simpleTeacherList) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

type VoteMeta struct {
	TeacherID int64 // 教师ID
	Score     int   // 评分
}

type ScoreInfo struct {
	votes   []int
	Average float64 `json:"average"`
}

// analyzer for choice or selection type
type TeacherScore struct {
	Average float64
	Total   int
	Count   int
	Meta    map[int]sourceList // option and its count
	Remark  []string           // remark for teacher
}

type sourceMeta Filter
type sourceList []sourceMeta

func (sm sourceList) Len() int {
	return len(sm)
}
func (sm sourceList) Swap(i, j int) {
	sm[i], sm[j] = sm[j], sm[i]
}
func (sm sourceList) Less(i, j int) bool {
	l := sm[i]
	r := sm[j]
	if l.Grade < r.Grade {
		return true
	} else if l.Grade > r.Grade {
		return false
	} else {
		return l.Index < r.Index
	}
}
