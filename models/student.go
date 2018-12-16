package models

import (
	"errors"
	"sort"
	"sync"

	"github.com/arong/dean/base"
	"github.com/astaxie/beego/logs"
)

// Um user manager
var Um userManager

type userManager struct {
	idMap   map[int64]*StudentInfo
	uuidMap map[string]*StudentInfo
	mutex   sync.Mutex
}

// Init: Init
func (um *userManager) Init(userMap map[int64]*StudentInfo) {
	um.uuidMap = make(map[string]*StudentInfo)
	if userMap != nil {
		um.idMap = userMap
		for _, v := range userMap {
			um.uuidMap[v.RegisterID] = v
		}
	} else {
		um.idMap = make(map[int64]*StudentInfo)
	}
}

// AddUser: AddUser
func (um *userManager) AddUser(u *StudentInfo) (int64, error) {
	var err error
	if len(u.RealName) == 0 {
		return 0, errors.New("invalid name")
	}

	if u.Gender < eGenderMale || u.Gender > eGenderUnknown {
		return 0, errors.New("invalid gender")
	}

	u.StudentID, err = Ma.InsertStudent(u)
	if err != nil {
		logs.Info("[AddUser]add user failed", err)
		return 0, err
	}

	um.idMap[u.StudentID] = u

	return u.StudentID, nil
}

// DelUser: DelUser
func (um *userManager) DelUser(uidList []int64) error {
	for _, uid := range uidList {
		_, ok := um.idMap[uid]
		if !ok {
			return errNotExist
		}

		err := Ma.DeleteStudent(uid)
		if err != nil {
			logs.Warn("[userManager::DelUser] failed", err)
			return err
		}
		delete(um.idMap, uid)
	}
	return nil
}

// ModUser: ModUser
func (um *userManager) ModUser(u *StudentInfo) error {
	if u.StudentID == 0 {
		return errors.New("invalid user id")
	}

	curr, ok := um.idMap[u.StudentID]
	if !ok {
		return errNotExist
	}

	if u.RegisterID != "" && curr.RegisterID != u.RegisterID {
		curr.RegisterID = u.RegisterID
	}
	return nil
}

// GetUser: GetUser
func (um *userManager) GetUser(uid int64) (*StudentInfo, error) {
	if val, ok := um.idMap[uid]; ok {
		return val, nil
	}
	return nil, errors.New("User not exists")
}

// GetUserByName: GetUserByName
func (um *userManager) GetUserByName(name string) (*StudentInfo, error) {
	return nil, errNotExist
}

func (um *userManager) GetStudentByRegisterNumber(reg string) (*StudentInfo, error) {
	s, ok := um.uuidMap[reg]
	if !ok {
		return nil, errNotExist
	}
	return s, nil
}

// GetAllUsers: GetAllUsers
func (um *userManager) GetAllUsers(f *StudentFilter) *base.CommList {
	resp := &base.CommList{}
	ret := studentList{}
	total := 0
	start, end := f.GetRange()

	curr := 0
	for _, v := range um.idMap {

		if f.Name != "" && f.Name != v.RealName {
			continue
		}

		if f.Number != "" && f.Number != v.RegisterID {
			continue
		}

		// handle page
		total++
		if curr < start || curr >= end {
			curr++
			continue
		}
		curr++
		ret = append(ret, v)
	}
	logs.Debug("[GetAllStudent]", "total count", len(um.idMap), "start", start, "end", end)
	sort.Sort(ret)
	resp.List = ret
	resp.Total = total
	return resp
}

func (um *userManager) getStudentList(grade int) ([]int64, error) {
	ret := []int64{}
	for _, v := range um.idMap {
		c, err := Cm.GetInfo(v.ClassID)
		if err != nil {
			continue
		}
		if grade != 0 && c.Grade != grade {
			continue
		}
		ret = append(ret, v.StudentID)
	}
	return ret, nil
}

// IsExist: IsExist
func (um *userManager) IsExist(studentID int64) bool {
	_, ok := um.idMap[studentID]
	return ok
}
