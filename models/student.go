package models

import (
	"errors"
	"github.com/arong/dean/base"
	"github.com/astaxie/beego/logs"
	"sort"
	"sync"
)

var Um userManager

type userManager struct {
	idMap map[int64]*StudentInfo
	mutex sync.Mutex
}

func (um *userManager) Init(userMap map[int64]*StudentInfo) {
	if userMap != nil {
		um.idMap = userMap
	} else {
		um.idMap = make(map[int64]*StudentInfo)
	}
}

type StudentInfo struct {
	profile
	ClassID    int    `json:"class_id"`
	StudentID  int64  `json:"student_id"`
	RegisterID string `json:"register_id"` // 学号
}

type userList []*StudentInfo

func (cl userList) Len() int {
	return len(cl)
}

func (cl userList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

func (cl userList) Less(i, j int) bool {
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

type UserID int64

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

func (um *userManager) DelUser(uidList []int64) error {
	for _, uid := range uidList {
		_, ok := um.idMap[uid]
		if !ok {
			return errNotExist
		}

		err := Ma.DeleteUser(uid)
		if err != nil {
			logs.Warn("[userManager::DelUser] failed", err)
			return err
		}
		delete(um.idMap, uid)
	}
	return nil
}

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

func (um *userManager) GetUser(uid int64) (*StudentInfo, error) {
	if val, ok := um.idMap[uid]; ok {
		return val, nil
	}
	return nil, errors.New("User not exists")
}

func (um *userManager) GetUserByName(name string) (*StudentInfo, error) {
	return nil, errNotExist
}

func (um *userManager) GetAllUsers(f *StudentFilter) *base.CommList {
	resp := &base.CommList{}
	ret := userList{}
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
