package models

import (
	"sort"
	"time"

	"github.com/bearbin/go-age"

	"github.com/astaxie/beego"

	"github.com/arong/dean/base"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
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
	logs.Info("[]", "newList", newList)
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

type Teacher struct {
	TeacherID int64  `json:"teacher_id"`
	SubjectID int    `json:"-"`
	Subject   string `json:"subject"`
	profile
}

func (t *Teacher) IsValid() error {
	if t.Mobile == "" {
		return errors.New("invalid mobile")
	}

	if t.Gender < eGenderMale && t.Gender > eGenderUnknown {
		return errors.New("invalid gender")
	}

	if t.Birthday == "" {
		return errors.New("empty birthday")
	}

	birth, err := time.Parse("2006-01-02", t.Birthday)
	if err != nil {
		return errors.New("invalid birthday")
	}
	t.Age = age.Age(birth)
	return nil
}

type TeacherList []*Teacher

func (tl TeacherList) Len() int {
	return len(tl)
}

func (tl TeacherList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (tl TeacherList) Less(i, j int) bool {
	return tl[i].TeacherID < tl[j].TeacherID
}

type TeacherFilter struct {
	base.CommPage
	Gender int    `json:"gender"`
	Age    int    `json:"age"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
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
